package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
)

type InstructoresController struct {
	beego.Controller
}

// GET /api/v1/instructores
func (c *InstructoresController) Listar() {
	var instructores []models.Instructor
	if !listarColeccion(&c.Controller, "instructores", &instructores) {
		return
	}

	responder(&c.Controller, 200, instructores)
}

// POST /api/v1/instructores
func (c *InstructoresController) Crear() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	var instructor models.Instructor

	if !leerCuerpo(&c.Controller, &instructor) {
		return
	}

	if strings.TrimSpace(instructor.Nombre) == "" {
		responderError(&c.Controller, 400, "El nombre es obligatorio.")
		return
	}

	crearEnColeccion(&c.Controller, "instructores", instructor)
}

// PUT /api/v1/instructores/:id
func (c *InstructoresController) Actualizar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	actualizarEnColeccion(&c.Controller, "instructores")
}

// DELETE /api/v1/instructores/:id
func (c *InstructoresController) Borrar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "instructores")
}
