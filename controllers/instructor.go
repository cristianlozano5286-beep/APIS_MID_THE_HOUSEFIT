package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type InstructorController struct {
	beego.Controller
	instructorService *services.InstructorService
}

func (c *InstructorController) Prepare() {
	c.instructorService = services.NewInstructorService()
}

// @Title Listar instructores
// @Description Lista instructores con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	gimnasio_id	query	int	false	"Filtrar por gimnasio"
// @Param	activo	query	bool	false	"Filtrar por estado activo"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *InstructorController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	gimnasioIDStr := c.GetString("gimnasio_id")
	activoStr := c.GetString("activo")

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	filters := map[string]interface{}{}
	if gimnasioIDStr != "" {
		gimnasioID, _ := strconv.Atoi(gimnasioIDStr)
		filters["gimnasio_id"] = gimnasioID
	}
	if activoStr != "" {
		activo := activoStr == "true"
		filters["activo"] = activo
	}

	instructores, total, err := c.instructorService.List(filters, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar instructores", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Instructores obtenidos", instructores, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}






