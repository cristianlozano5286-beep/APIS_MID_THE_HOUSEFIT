package controllers

import (
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
