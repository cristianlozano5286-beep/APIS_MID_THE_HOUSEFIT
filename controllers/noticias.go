package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type NoticiaController struct {
	beego.Controller
	noticiaService *services.NoticiaService
}

func (c *NoticiaController) Prepare() {
	c.noticiaService = services.NewNoticiaService()
}

// @Title Listar noticias
// @Description Lista noticias con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	tipo	query	string	false	"Filtrar por tipo"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *NoticiaController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	tipo := c.GetString("tipo")

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	// Los usuarios no administrador ven todas las noticias (no hay campo "publicada" en BD)
	filters := map[string]interface{}{}
	if tipo != "" {
		filters["tipo"] = tipo
	}

	noticias, total, err := c.noticiaService.List(filters, pagina, porPagina, false)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar noticias", err)
		c.ServeJSON()
		return
	}