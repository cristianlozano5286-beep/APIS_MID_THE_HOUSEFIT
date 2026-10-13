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
