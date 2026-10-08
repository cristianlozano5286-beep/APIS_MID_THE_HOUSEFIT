package controllers

import (
	"net/http"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type DashboardController struct {
	beego.Controller
	dashboardService *services.DashboardService
}

func (c *DashboardController) Prepare() {
	c.dashboardService = services.NewDashboardService()
}

// @Title Métricas del dashboard
// @Description Obtiene las métricas principales para el dashboard (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /metricas [get]
func (c *DashboardController) Metricas() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	metricas, err := c.dashboardService.GetMetricas()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener métricas", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Métricas obtenidas", metricas)
	c.Data["json"] = response
	c.ServeJSON()
}