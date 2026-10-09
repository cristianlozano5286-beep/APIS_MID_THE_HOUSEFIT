package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type GimnasioController struct {
	beego.Controller
	gimnasioService *services.GimnasioService
}

func (c *GimnasioController) Prepare() {
	c.gimnasioService = services.NewGimnasioService()
}

// @Title Listar gimnasios
// @Description Lista gimnasios con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	ciudad	query	string	false	"Filtrar por ciudad"
// @Param	activo	query	bool	false	"Filtrar por estado activo"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *GimnasioController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	ciudad := c.GetString("ciudad")
	activoStr := c.GetString("activo")

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	filters := map[string]interface{}{}
	if ciudad != "" {
		filters["ciudad"] = ciudad
	}
	if activoStr != "" {
		activo := activoStr == "true"
		filters["activo"] = activo
	}

	gimnasios, total, err := c.gimnasioService.List(filters, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar gimnasios", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Gimnasios obtenidos", gimnasios, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Ver gimnasio
// @Description Obtiene un gimnasio por ID con fotos y servicios
// @Param	id	path	int	true	"ID del gimnasio"
// @Success 200 {object} models.RespuestaAPI
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [get]
func (c *GimnasioController) Ver() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	gimnasio, err := c.gimnasioService.GetCompleto(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener gimnasio", err)
		c.ServeJSON()
		return
	}

	if gimnasio == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Gimnasio no encontrado", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Gimnasio obtenido", gimnasio)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Crear gimnasio
// @Description Crea un nuevo gimnasio (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.GimnasioRequest	true	"Datos del gimnasio"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [post]
func (c *GimnasioController) Crear() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	var req models.GimnasioRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	gimnasio, err := c.gimnasioService.Create(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al crear gimnasio", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Gimnasio creado correctamente", gimnasio)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Actualizar gimnasio
// @Description Actualiza un gimnasio (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID del gimnasio"
// @Param	body	body 	models.GimnasioUpdateRequest	true	"Datos a actualizar"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [put]
func (c *GimnasioController) Actualizar() {
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

	var req models.GimnasioUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	gimnasio, err := c.gimnasioService.Update(id, req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al actualizar gimnasio", err)
		c.ServeJSON()
		return
	}

	if gimnasio == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Gimnasio no encontrado", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Gimnasio actualizado correctamente", gimnasio)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar gimnasio
// @Description Elimina un gimnasio (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID del gimnasio"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *GimnasioController) Eliminar() {
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

	gimnasio, err := c.gimnasioService.GetCompleto(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar gimnasio", err)
		c.ServeJSON()
		return
	}
	if gimnasio == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Gimnasio no encontrado", nil)
		c.ServeJSON()
		return
	}

	err = c.gimnasioService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar gimnasio", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Gimnasio eliminado correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Obtener clases de un gimnasio
// @Description Obtiene las clases disponibles de un gimnasio
// @Param	id	path	int	true	"ID del gimnasio"
// @Success 200 {object} models.RespuestaAPI
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id/clases [get]
func (c *GimnasioController) ObtenerClases() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	claseService := services.NewClaseService()
	clases, err := claseService.GetByGimnasio(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener clases", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Clases obtenidas", clases)
	c.Data["json"] = response
	c.ServeJSON()
}