package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/go-playground/validator/v10"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type RutinaController struct {
	beego.Controller
	rutinaService *services.RutinaService
}

func (c *RutinaController) Prepare() {
	c.rutinaService = services.NewRutinaService()
}

// @Title Listar rutinas
// @Description Lista rutinas con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	nivel	query	string	false	"Filtrar por nivel"
// @Param	tipo_entrenamiento	query	string	false	"Filtrar por tipo de entrenamiento"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *RutinaController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	nivel := c.GetString("nivel")
	tipoEntrenamiento := c.GetString("tipo_entrenamiento")

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	filters := map[string]interface{}{}
	if nivel != "" {
		filters["nivel"] = nivel
	}
	if tipoEntrenamiento != "" {
		filters["tipo_entrenamiento"] = tipoEntrenamiento
	}

	rutinas, total, err := c.rutinaService.List(filters, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar rutinas", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Rutinas obtenidas", rutinas, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Ver rutina
// @Description Obtiene una rutina por ID con sus ejercicios
// @Param	id	path	int	true	"ID de la rutina"
// @Success 200 {object} models.RespuestaAPI
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [get]
func (c *RutinaController) Ver() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	rutina, err := c.rutinaService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener rutina", err)
		c.ServeJSON()
		return
	}

	if rutina == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Rutina no encontrada", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Rutina obtenida", rutina)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Crear rutina
// @Description Crea una nueva rutina con ejercicios (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.RutinaRequest	true	"Datos de la rutina"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [post]
func (c *RutinaController) Crear() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	var req models.RutinaRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	rutina, err := c.rutinaService.Create(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al crear rutina", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Rutina creada correctamente", rutina)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Actualizar rutina
// @Description Actualiza una rutina (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la rutina"
// @Param	body	body 	models.RutinaUpdateRequest	true	"Datos a actualizar"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [put]
func (c *RutinaController) Actualizar() {
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

	var req models.RutinaUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	validate := validator.New()
    if err := validate.Struct(&req); err != nil {
        c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
        c.Data["json"] = models.RespuestaErrorGeneral("Los datos no cumplen con las reglas de validación", err)
        c.ServeJSON()
        return
    }

	rutina, err := c.rutinaService.Update(id, req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al actualizar rutina", err)
		c.ServeJSON()
		return
	}

	if rutina == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Rutina no encontrada", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Rutina actualizada correctamente", rutina)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar rutina
// @Description Elimina una rutina y sus ejercicios (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la rutina"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *RutinaController) Eliminar() {
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

	rutina, err := c.rutinaService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar rutina", err)
		c.ServeJSON()
		return
	}
	if rutina == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Rutina no encontrada", nil)
		c.ServeJSON()
		return
	}

	err = c.rutinaService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar rutina", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Rutina eliminada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}