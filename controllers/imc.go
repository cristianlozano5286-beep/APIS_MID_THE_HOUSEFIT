package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type IMCController struct {
	beego.Controller
	imcService *services.IMCService
}

func (c *IMCController) Prepare() {
	c.imcService = services.NewIMCService()
}

// @Title Calcular IMC
// @Description Calcula el IMC y devuelve recomendaciones
// @Param	body	body 	models.CalcularIMCRequest	true	"Datos para cálculo"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /calcular [post]
func (c *IMCController) Calcular() {
	var req models.CalcularIMCRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	if req.UsuarioCorreo == nil {
		usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
		if usuarioCorreo != nil {
			correo := usuarioCorreo.(string)
			req.UsuarioCorreo = &correo
		}
	}

	resultado, err := c.imcService.CalcularIMC(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al calcular IMC", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("IMC calculado correctamente", resultado)
	c.Data["json"] = response
	c.ServeJSON()
}


// @Title Historial IMC
// @Description Obtiene el historial de IMC del usuario autenticado
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 401 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /historial [get]
func (c *IMCController) Historial() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
	if usuarioCorreo == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
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

	resultados, total, err := c.imcService.GetHistorial(usuarioCorreo.(string), pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener historial IMC", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Historial IMC obtenido", resultados, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}








