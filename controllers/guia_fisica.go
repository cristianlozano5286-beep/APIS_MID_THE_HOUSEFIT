package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type GuiaFisicaController struct {
	beego.Controller
	guiaService *services.GuiaFisicaService
}

func (c *GuiaFisicaController) Prepare() {
	c.guiaService = services.NewGuiaFisicaService()
}

// @Title Listar guías físicas
// @Description Lista guías físicas con paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *GuiaFisicaController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	guias, total, err := c.guiaService.List(pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar guías físicas", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Guías físicas obtenidas", guias, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Ver guía física
// @Description Obtiene una guía física por ID
// @Param	id	path	int	true	"ID de la guía"
// @Success 200 {object} models.RespuestaAPI
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [get]
func (c *GuiaFisicaController) Ver() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	guia, err := c.guiaService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener guía física", err)
		c.ServeJSON()
		return
	}

	if guia == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Guía física no encontrada", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Guía física obtenida", guia)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Crear guía física
// @Description Crea una nueva guía física (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.GuiaFisicaRequest	true	"Datos de la guía"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [post]
func (c *GuiaFisicaController) Crear() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	var req models.GuiaFisicaRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	guia, err := c.guiaService.Create(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al crear guía física", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Guía física creada correctamente", guia)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Actualizar guía física
// @Description Actualiza una guía física (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la guía"
// @Param	body	body 	models.GuiaFisicaUpdateRequest	true	"Datos a actualizar"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [put]
func (c *GuiaFisicaController) Actualizar() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	var req models.GuiaFisicaUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	guia, err := c.guiaService.Update(id, req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al actualizar guía física", err)
		c.ServeJSON()
		return
	}

	if guia == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Guía física no encontrada", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Guía física actualizada correctamente", guia)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar guía física
// @Description Elimina una guía física (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la guía"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *GuiaFisicaController) Eliminar() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	guia, err := c.guiaService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar guía física", err)
		c.ServeJSON()
		return
	}
	if guia == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Guía física no encontrada", nil)
		c.ServeJSON()
		return
	}

	err = c.guiaService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar guía física", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Guía física eliminada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}