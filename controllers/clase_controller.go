package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type ClaseController struct {
	beego.Controller
	claseService   *services.ClaseService
	reservaService *services.ReservaService
}

func (c *ClaseController) Prepare() {
	c.claseService = services.NewClaseService()
	c.reservaService = services.NewReservaService()
}

// @Title Listar clases
// @Description Lista clases disponibles con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	gimnasio_id	query	int	false	"Filtrar por gimnasio"
// @Param	tipo_clase	query	string	false	"Filtrar por tipo de clase"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *ClaseController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	gimnasioIDStr := c.GetString("gimnasio_id")
	tipoClase := c.GetString("tipo_clase")

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
	if tipoClase != "" {
		filters["tipo_clase"] = tipoClase
	}

	clases, total, err := c.claseService.List(filters, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar clases", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Clases obtenidas", clases, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}