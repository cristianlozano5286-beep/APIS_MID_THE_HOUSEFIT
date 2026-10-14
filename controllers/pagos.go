package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type PagoController struct {
	beego.Controller
	pagoService *services.PagoService
}

func (c *PagoController) Prepare() {
	c.pagoService = services.NewPagoService()
}

// @Title Listar métodos de pago
// @Description Lista los métodos de pago disponibles
// @Success 200 {object} models.RespuestaAPI
// @Failure 500 {object} models.RespuestaError
// @router /metodos [get]
func (c *PagoController) ListarMetodos() {
	metodos, err := c.pagoService.ListarMetodos()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar métodos de pago", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Métodos de pago obtenidos", metodos)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Procesar pago
// @Description Procesa un pago (simulado)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.ProcesarPagoRequest	true	"Datos del pago"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /procesar [post]
func (c *PagoController) ProcesarPago() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
	if usuarioCorreo == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
		c.ServeJSON()
		return
	}

	var req models.ProcesarPagoRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	req.UsuarioCorreo = usuarioCorreo.(string)

	resultado, err := c.pagoService.ProcesarPago(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al procesar pago", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Pago procesado", resultado)
	c.Data["json"] = response
	c.ServeJSON()
}