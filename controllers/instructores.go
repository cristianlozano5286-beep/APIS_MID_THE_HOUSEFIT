package controllers

import (
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
