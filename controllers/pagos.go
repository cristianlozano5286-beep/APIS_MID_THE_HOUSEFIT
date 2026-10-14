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