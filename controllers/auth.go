package controllers

import (
	"crypto/rand"
	"fmt"
	"html"
	"math/big"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"

	"thehousefit_mid/models"
	"thehousefit_mid/services"
)

// los codigos se guardan en la ram, no en la base
var codigos = make(map[string]models.CodigoVerificacion)

// candado para que no se crucen las peticiones
var candadoCodigos sync.RWMutex

const maximoIntentos = 5

type AuthController struct {
	beego.Controller
}

// POST /api/v1/auth/login
func (c *AuthController) Login() {
	var solicitud models.SolicitudLogin

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := limpiarCorreo(solicitud.Correo)

	usuario, err := buscarUsuarioPorCorreo(correo)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// el mismo mensaje si el correo no existe o la contraseña esta mal
	if usuario == nil || bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(solicitud.Password)) != nil {
		responderError(&c.Controller, 401, "Correo o contraseña incorrectos.")
		return
	}

	if !usuario.CorreoVerificado {
		responder(&c.Controller, 403, map[string]interface{}{
			"ok":                   false,
			"mensaje":              "Primero verifica tu correo.",
			"requiereVerificacion": true,
		})
		return
	}

	respuesta, err := armarSesion(usuario)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo crear el token.")
		return
	}

	responder(&c.Controller, 200, respuesta)
}

// POST /api/v1/auth/registro
func (c *AuthController) Registro() {
	var solicitud models.SolicitudRegistro

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	nombre := strings.TrimSpace(solicitud.Nombre)
	correo := limpiarCorreo(solicitud.Correo)

	if nombre == "" {
		responderError(&c.Controller, 400, "El nombre es obligatorio.")
		return
	}

	_, err := mail.ParseAddress(correo)
	if err != nil {
		responderError(&c.Controller, 400, "El correo no tiene un formato válido.")
		return
	}

	if len(solicitud.Password) < 6 {
		responderError(&c.Controller, 400, "La contraseña debe tener mínimo 6 caracteres.")
		return
	}

	existente, err := buscarUsuarioPorCorreo(correo)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	if existente != nil {
		responderError(&c.Controller, 409, "Ese correo ya está registrado.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(solicitud.Password), bcrypt.DefaultCost)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo proteger la contraseña.")
		return
	}

	// toda cuenta nueva empieza con el rol Usuario
	nuevo := models.Usuario{
		Nombre:           nombre,
		Correo:           correo,
		Password:         string(hash),
		Rol:              "Usuario",
		CorreoVerificado: false,
	}

	err = services.NuevoServicioCrud().Crear("usuarios", nuevo, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if !enviarCodigo(&c.Controller, correo, nombre) {
		return
	}

	responder(&c.Controller, 201, map[string]interface{}{
		"ok":      true,
		"mensaje": "Cuenta creada. Te enviamos un código al correo.",
	})
}

// POST /api/v1/auth/codigo
// sirve para pedir otro codigo de verificacion
func (c *AuthController) PedirCodigo() {
	var solicitud models.SolicitudCorreo

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := limpiarCorreo(solicitud.Correo)

	usuario, err := buscarUsuarioPorCorreo(correo)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	if usuario == nil {
		responderError(&c.Controller, 404, "No existe una cuenta con ese correo.")
		return
	}

	if !enviarCodigo(&c.Controller, correo, usuario.Nombre) {
		return
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"ok":      true,
		"mensaje": "Te enviamos un código al correo.",
	})
}

// POST /api/v1/auth/verificar
func (c *AuthController) VerificarCodigo() {
	var solicitud models.SolicitudVerificar

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := limpiarCorreo(solicitud.Correo)

	if !comprobarCodigo(&c.Controller, correo, solicitud.Codigo) {
		return
	}

	usuario, err := buscarUsuarioPorCorreo(correo)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	if usuario == nil {
		responderError(&c.Controller, 404, "No existe una cuenta con ese correo.")
		return
	}

	cambios := map[string]interface{}{"correoVerificado": true}
	err = services.NuevoServicioCrud().Actualizar("usuarios", strconv.Itoa(usuario.Id), cambios, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	usuario.CorreoVerificado = true

	respuesta, err := armarSesion(usuario)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo crear el token.")
		return
	}

	responder(&c.Controller, 200, respuesta)
}

// POST /api/v1/auth/recuperar
func (c *AuthController) Recuperar() {
	var solicitud models.SolicitudCorreo

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := limpiarCorreo(solicitud.Correo)

	usuario, err := buscarUsuarioPorCorreo(correo)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	if usuario == nil {
		responderError(&c.Controller, 404, "No existe una cuenta con ese correo.")
		return
	}

	if !enviarCodigo(&c.Controller, correo, usuario.Nombre) {
		return
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"ok":      true,
		"mensaje": "Te enviamos un código para cambiar tu contraseña.",
	})
}

// POST /api/v1/auth/restablecer
func (c *AuthController) Restablecer() {
	var solicitud models.SolicitudRestablecer

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := limpiarCorreo(solicitud.Correo)

	if len(solicitud.NuevaPassword) < 6 {
		responderError(&c.Controller, 400, "La contraseña debe tener mínimo 6 caracteres.")
		return
	}

	if !comprobarCodigo(&c.Controller, correo, solicitud.Codigo) {
		return
	}

	usuario, err := buscarUsuarioPorCorreo(correo)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	if usuario == nil {
		responderError(&c.Controller, 404, "No existe una cuenta con ese correo.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(solicitud.NuevaPassword), bcrypt.DefaultCost)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo proteger la contraseña.")
		return
	}

	// si pudo usar el codigo, el correo es suyo
	cambios := map[string]interface{}{
		"password":         string(hash),
		"correoVerificado": true,
	}
	err = services.NuevoServicioCrud().Actualizar("usuarios", strconv.Itoa(usuario.Id), cambios, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"ok":      true,
		"mensaje": "Tu contraseña fue cambiada.",
	})
}

// GET /api/v1/auth/perfil
func (c *AuthController) Perfil() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	var usuario models.Usuario
	err := services.NuevoServicioCrud().Obtener("usuarios", sesion.ID, &usuario)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, models.UsuarioAuth{
		Nombre: usuario.Nombre,
		Correo: usuario.Correo,
		Rol:    usuario.Rol,
	})
}

// ---------- funciones de apoyo ----------

func limpiarCorreo(correo string) string {
	return strings.ToLower(strings.TrimSpace(correo))
}

// busca en el CRUD al usuario con ese correo, devuelve nil si no existe
func buscarUsuarioPorCorreo(correo string) (*models.Usuario, error) {
	var usuarios []models.Usuario

	err := services.NuevoServicioCrud().Listar("usuarios", &usuarios)
	if err != nil {
		return nil, err
	}

	for i := range usuarios {
		if strings.ToLower(usuarios[i].Correo) == correo {
			return &usuarios[i], nil
		}
	}

	return nil, nil
}

// arma la respuesta del login: el token y los datos de la sesion
func armarSesion(usuario *models.Usuario) (map[string]interface{}, error) {
	token, err := services.GenerarToken(strconv.Itoa(usuario.Id), usuario.Correo, usuario.Nombre, usuario.Rol)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"usuario": models.UsuarioAuth{
			Nombre: usuario.Nombre,
			Correo: usuario.Correo,
			Rol:    usuario.Rol,
		},
	}, nil
}

// crea un codigo de 6 digitos y lo manda al correo
// si algo falla ya responde el error y devuelve false
func enviarCodigo(c *beego.Controller, correo string, nombre string) bool {
	// no se puede pedir otro codigo antes de 30 segundos
	candadoCodigos.RLock()
	anterior, hayAnterior := codigos[correo]
	candadoCodigos.RUnlock()

	if hayAnterior && time.Since(anterior.Creado) < 30*time.Second {
		responderError(c, 429, "Espera unos segundos para pedir otro código.")
		return false
	}

	codigo, err := generarCodigo()
	if err != nil {
		responderError(c, 500, "No se pudo generar el código.")
		return false
	}

	candadoCodigos.Lock()
	codigos[correo] = models.CodigoVerificacion{
		Correo:   correo,
		Codigo:   codigo,
		Creado:   time.Now(),
		Vence:    time.Now().Add(10 * time.Minute),
		Intentos: 0,
		Activo:   true,
	}
	candadoCodigos.Unlock()

	err = mandarCodigo(correo, codigo, nombre)
	if err != nil {
		candadoCodigos.Lock()
		delete(codigos, correo)
		candadoCodigos.Unlock()

		responderError(c, 500, "No se pudo enviar el correo: "+err.Error())
		return false
	}

	return true
}

// revisa el codigo que escribio la persona
// si algo esta mal ya responde el error y devuelve false
func comprobarCodigo(c *beego.Controller, correo string, codigo string) bool {
	codigo = strings.TrimSpace(codigo)

	if correo == "" || codigo == "" {
		responderError(c, 400, "El correo y el código son obligatorios.")
		return false
	}

	candadoCodigos.RLock()
	guardado, existe := codigos[correo]
	candadoCodigos.RUnlock()

	if !existe {
		responderError(c, 404, "No encontramos un código para ese correo. Pídelo de nuevo.")
		return false
	}

	if !guardado.Activo {
		responderError(c, 401, "Ese código ya no sirve. Pide uno nuevo.")
		return false
	}

	if time.Now().After(guardado.Vence) {
		candadoCodigos.Lock()
		guardado.Activo = false
		codigos[correo] = guardado
		candadoCodigos.Unlock()

		responderError(c, 401, "El código se venció. Pide uno nuevo.")
		return false
	}

	if codigo != guardado.Codigo {
		candadoCodigos.Lock()
		guardado.Intentos++
		if guardado.Intentos >= maximoIntentos {
			guardado.Activo = false
		}
		codigos[correo] = guardado
		candadoCodigos.Unlock()

		if !guardado.Activo {
			responderError(c, 401, "Demasiados intentos. Pide un código nuevo.")
			return false
		}

		quedan := maximoIntentos - guardado.Intentos
		responderError(c, 401, fmt.Sprintf("El código no coincide. Te quedan %d intentos.", quedan))
		return false
	}

	// codigo bueno, ya no se puede volver a usar
	candadoCodigos.Lock()
	guardado.Activo = false
	codigos[correo] = guardado
	candadoCodigos.Unlock()

	return true
}

func generarCodigo() (string, error) {
	numero, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", numero.Int64()+100000), nil
}

func mandarCodigo(correo string, codigo string, nombre string) error {
	if !beego.AppConfig.DefaultBool("enviar_correo", false) {
		fmt.Println("Código para", correo, "->", codigo)
		return nil
	}

	servicioGmail, err := services.NuevoServicioGmail()
	if err != nil {
		return err
	}

	plantilla, err := os.ReadFile("templates/correo_codigo.html")
	if err != nil {
		return err
	}

	saludo := "Hola,"
	if strings.TrimSpace(nombre) != "" {
		saludo = "Hola, " + html.EscapeString(strings.Fields(nombre)[0]) + ","
	}

	cuerpo := strings.ReplaceAll(string(plantilla), "{{CODIGO}}", codigo)
	cuerpo = strings.ReplaceAll(cuerpo, "{{SALUDO}}", saludo)

	// el correo no lleva logo, por eso va vacio
	_, err = servicioGmail.EnviarCorreo(correo, "Tu código de The House Fit", cuerpo, nil)
	return err
}
