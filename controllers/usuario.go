package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
	authService *services.AuthService
}

func (c *UserController) Prepare() {
	c.authService = services.NewAuthService()
}

// @Title Listar usuarios
// @Description Lista todos los usuarios (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @router / [get]

func (c *UserController) Listar() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	usuarios, total, err := c.authService.ListarUsuarios(pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar usuarios", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Usuarios obtenidos", usuarios, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}
