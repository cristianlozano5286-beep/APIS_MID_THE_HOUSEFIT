package controllers

import (
	"math"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
)

// rutinas, guias fisicas, guias nutricionales y la calculadora de IMC
type ContenidoController struct {
	beego.Controller
}

// GET /api/v1/rutinas?tipo=Fuerza
func (c *ContenidoController) ListarRutinas() {
	var rutinas []models.Rutina
	if !listarColeccion(&c.Controller, "rutinas", &rutinas) {
		return
	}

	tipo := strings.TrimSpace(c.GetString("tipo"))

	resultado := []models.Rutina{}
	for _, r := range rutinas {
		if tipo != "" && r.TipoEntrenamiento != tipo {
			continue
		}
		resultado = append(resultado, r)
	}

	responder(&c.Controller, 200, resultado)
}

// POST /api/v1/rutinas
func (c *ContenidoController) CrearRutina() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	var rutina models.Rutina

	if !leerCuerpo(&c.Controller, &rutina) {
		return
	}

	if strings.TrimSpace(rutina.Nombre) == "" {
		responderError(&c.Controller, 400, "El nombre es obligatorio.")
		return
	}

	crearEnColeccion(&c.Controller, "rutinas", rutina)
}

// DELETE /api/v1/rutinas/:id
func (c *ContenidoController) BorrarRutina() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "rutinas")
}

// GET /api/v1/guias-fisicas
func (c *ContenidoController) ListarGuiasFisicas() {
	var guias []models.GuiaFisica
	if !listarColeccion(&c.Controller, "guias-fisicas", &guias) {
		return
	}

	responder(&c.Controller, 200, guias)
}

// POST /api/v1/guias-fisicas
func (c *ContenidoController) CrearGuiaFisica() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	var guia models.GuiaFisica

	if !leerCuerpo(&c.Controller, &guia) {
		return
	}

	if strings.TrimSpace(guia.Titulo) == "" || strings.TrimSpace(guia.Contenido) == "" {
		responderError(&c.Controller, 400, "El título y el contenido son obligatorios.")
		return
	}

	guia.FechaPublicacion = hoy()

	crearEnColeccion(&c.Controller, "guias-fisicas", guia)
}

// DELETE /api/v1/guias-fisicas/:id
func (c *ContenidoController) BorrarGuiaFisica() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "guias-fisicas")
}

// GET /api/v1/guias-nutricionales
func (c *ContenidoController) ListarGuiasNutricionales() {
	var guias []models.GuiaNutricional
	if !listarColeccion(&c.Controller, "guias-nutricionales", &guias) {
		return
	}

	responder(&c.Controller, 200, guias)
}

// POST /api/v1/guias-nutricionales
func (c *ContenidoController) CrearGuiaNutricional() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	var guia models.GuiaNutricional

	if !leerCuerpo(&c.Controller, &guia) {
		return
	}

	if strings.TrimSpace(guia.Titulo) == "" || strings.TrimSpace(guia.Recomendaciones) == "" {
		responderError(&c.Controller, 400, "El título y las recomendaciones son obligatorios.")
		return
	}

	crearEnColeccion(&c.Controller, "guias-nutricionales", guia)
}

// DELETE /api/v1/guias-nutricionales/:id
func (c *ContenidoController) BorrarGuiaNutricional() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "guias-nutricionales")
}

// GET /api/v1/imc?peso=70&estatura=1.75
// calcula el IMC y devuelve la guia nutricional que le toca
func (c *ContenidoController) CalcularIMC() {
	peso, errPeso := strconv.ParseFloat(c.GetString("peso"), 64)
	estatura, errEstatura := strconv.ParseFloat(c.GetString("estatura"), 64)

	if errPeso != nil || errEstatura != nil || peso <= 0 || estatura <= 0 {
		responderError(&c.Controller, 400, "Escribe el peso en kg y la estatura en metros.")
		return
	}

	imc := math.Round(peso/(estatura*estatura)*10) / 10
	categoria := clasificarIMC(imc)

	var guias []models.GuiaNutricional
	if !listarColeccion(&c.Controller, "guias-nutricionales", &guias) {
		return
	}

	// puede que todavia no haya guia para esa categoria
	var guia *models.GuiaNutricional
	for i := range guias {
		if guias[i].CategoriaIMC == categoria {
			guia = &guias[i]
		}
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"imc":       imc,
		"categoria": categoria,
		"guia":      guia,
	})
}

func clasificarIMC(imc float64) string {
	if imc < 18.5 {
		return "Bajo peso"
	}
	if imc < 25 {
		return "Peso normal"
	}
	if imc < 30 {
		return "Sobrepeso"
	}
	return "Obesidad"
}
