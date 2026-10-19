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
