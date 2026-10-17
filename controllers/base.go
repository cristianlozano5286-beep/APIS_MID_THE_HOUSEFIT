package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/services"
)

func responder(c *beego.Controller, codigo int, datos interface{}) {
	c.Ctx.Output.SetStatus(codigo)
	c.Data["json"] = datos
	c.ServeJSON()
}

func responderError(c *beego.Controller, codigo int, mensaje string) {
	responder(c, codigo, map[string]interface{}{
		"ok":      false,
		"mensaje": mensaje,
	})
}

// para las respuestas sin cuerpo
func responderVacio(c *beego.Controller) {
	c.Ctx.Output.SetStatus(204)
	c.Ctx.Output.Body([]byte(""))
}

// responde segun lo que pasó al hablar con el CRUD
func errorCrud(c *beego.Controller, err error) {
	if errors.Is(err, services.ErrNoEncontrado) {
		responderError(c, 404, "No se encontró lo que buscas.")
		return
	}
	if errors.Is(err, services.ErrDatoInvalido) {
		responderError(c, 400, "Los datos enviados no son válidos.")
		return
	}
	if errors.Is(err, services.ErrYaExiste) {
		responderError(c, 409, "Ya existe un registro igual.")
		return
	}
	responderError(c, 502, "El servicio de datos no respondió: "+err.Error())
}

// lee el JSON del cuerpo, si viene mal ya responde el error
func leerCuerpo(c *beego.Controller, destino interface{}) bool {
	cuerpo := c.Ctx.Input.RequestBody

	if len(cuerpo) == 0 {
		responderError(c, 400, "El cuerpo de la solicitud está vacío.")
		return false
	}

	err := json.Unmarshal(cuerpo, destino)
	if err != nil {
		responderError(c, 400, "El cuerpo de la solicitud no tiene un formato válido.")
		return false
	}

	return true
}

func hoy() string {
	return time.Now().Format("2006-01-02")
}

// el id de la ruta tiene que ser un numero
func leerId(c *beego.Controller) (string, bool) {
	id := c.Ctx.Input.Param(":id")

	_, err := strconv.Atoi(id)
	if err != nil {
		responderError(c, 400, "El id debe ser un número.")
		return "", false
	}

	return id, true
}

// ---------- tokens ----------

// mira el token pero no responde nada si falla
func leerToken(c *beego.Controller) (*services.Claims, bool) {
	cabecera := c.Ctx.Input.Header("Authorization")

	if !strings.HasPrefix(cabecera, "Bearer ") {
		return nil, false
	}

	claims, err := services.ValidarToken(strings.TrimPrefix(cabecera, "Bearer "))
	if err != nil {
		return nil, false
	}

	return claims, true
}

// deja pasar a cualquiera que haya iniciado sesion
func exigirSesion(c *beego.Controller) (*services.Claims, bool) {
	sesion, valido := leerToken(c)

	if !valido {
		responderError(c, 401, "Inicia sesión para continuar.")
		return nil, false
	}

	return sesion, true
}

// deja pasar solo a los roles que se le indiquen
// ejemplo: exigirRol(&c.Controller, "Administrador", "Entrenador")
func exigirRol(c *beego.Controller, roles ...string) (*services.Claims, bool) {
	sesion, valido := exigirSesion(c)
	if !valido {
		return nil, false
	}

	for _, rol := range roles {
		if sesion.Rol == rol {
			return sesion, true
		}
	}

	responderError(c, 403, "No tienes permiso para hacer esto.")
	return nil, false
}

// ---------- funciones que se repiten en casi todos los controladores ----------

// trae toda una coleccion del CRUD y la deja en destino
func listarColeccion(c *beego.Controller, coleccion string, destino interface{}) bool {
	err := services.NuevoServicioCrud().Listar(coleccion, destino)
	if err != nil {
		errorCrud(c, err)
		return false
	}
	return true
}

// GET /coleccion/:id
func obtenerDeColeccion(c *beego.Controller, coleccion string) {
	id, valido := leerId(c)
	if !valido {
		return
	}

	var encontrado map[string]interface{}
	err := services.NuevoServicioCrud().Obtener(coleccion, id, &encontrado)
	if err != nil {
		errorCrud(c, err)
		return
	}

	responder(c, 200, encontrado)
}

// POST /coleccion, el dato ya viene leido y revisado
func crearEnColeccion(c *beego.Controller, coleccion string, dato interface{}) {
	var creado map[string]interface{}
	err := services.NuevoServicioCrud().Crear(coleccion, dato, &creado)
	if err != nil {
		errorCrud(c, err)
		return
	}

	responder(c, 201, creado)
}

// PUT /coleccion/:id, solo se cambia lo que venga en el cuerpo
func actualizarEnColeccion(c *beego.Controller, coleccion string) {
	id, valido := leerId(c)
	if !valido {
		return
	}

	var cambios map[string]interface{}
	if !leerCuerpo(c, &cambios) {
		return
	}

	// el id no se puede cambiar
	delete(cambios, "id")

	var actualizado map[string]interface{}
	err := services.NuevoServicioCrud().Actualizar(coleccion, id, cambios, &actualizado)
	if err != nil {
		errorCrud(c, err)
		return
	}

	responder(c, 200, actualizado)
}

// DELETE /coleccion/:id
func borrarDeColeccion(c *beego.Controller, coleccion string) {
	id, valido := leerId(c)
	if !valido {
		return
	}

	err := services.NuevoServicioCrud().Borrar(coleccion, id)
	if err != nil {
		errorCrud(c, err)
		return
	}

	responderVacio(c)
}
