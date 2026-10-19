package controllers

import (
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
)

type ClasesController struct {
	beego.Controller
}

// GET /api/v1/clases?gimnasioId=1&tipoClase=Yoga
func (c *ClasesController) Listar() {
	var clases []models.Clase
	if !listarColeccion(&c.Controller, "clases", &clases) {
		return
	}

	gimnasioId := c.GetString("gimnasioId")
	tipoClase := strings.TrimSpace(c.GetString("tipoClase"))

	resultado := []models.Clase{}
	for _, clase := range clases {
		if gimnasioId != "" && strconv.Itoa(clase.GimnasioId) != gimnasioId {
			continue
		}
		if tipoClase != "" && clase.TipoClase != tipoClase {
			continue
		}
		resultado = append(resultado, clase)
	}

	responder(&c.Controller, 200, resultado)
}

// POST /api/v1/clases
func (c *ClasesController) Crear() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	var clase models.Clase

	if !leerCuerpo(&c.Controller, &clase) {
		return
	}

	if clase.GimnasioId == 0 || strings.TrimSpace(clase.Instructor) == "" || strings.TrimSpace(clase.Hora) == "" {
		responderError(&c.Controller, 400, "El gimnasio, el instructor y la hora son obligatorios.")
		return
	}

	crearEnColeccion(&c.Controller, "clases", clase)
}

// DELETE /api/v1/clases/:id
func (c *ClasesController) Borrar() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "clases")
}
