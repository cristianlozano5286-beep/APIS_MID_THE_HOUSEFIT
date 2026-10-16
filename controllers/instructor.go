package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type InstructorController struct {
	beego.Controller
	instructorService *services.InstructorService
}

func (c *InstructorController) Prepare() {
	c.instructorService = services.NewInstructorService()
}

// @Title Listar instructores
// @Description Lista instructores con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	gimnasio_id	query	int	false	"Filtrar por gimnasio"
// @Param	activo	query	bool	false	"Filtrar por estado activo"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *InstructorController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	gimnasioIDStr := c.GetString("gimnasio_id")
	activoStr := c.GetString("activo")

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}

	filters := map[string]interface{}{}
	if gimnasioIDStr != "" {
		gimnasioID, _ := strconv.Atoi(gimnasioIDStr)
		filters["gimnasio_id"] = gimnasioID
	}
	if activoStr != "" {
		activo := activoStr == "true"
		filters["activo"] = activo
	}

	instructores, total, err := c.instructorService.List(filters, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar instructores", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Instructores obtenidos", instructores, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Ver instructor
// @Description Obtiene un instructor por ID
// @Param	id	path	int	true	"ID del instructor"
// @Success 200 {object} models.RespuestaAPI
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [get]
func (c *InstructorController) Ver() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	instructor, err := c.instructorService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener instructor", err)
		c.ServeJSON()
		return
	}

	if instructor == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Instructor no encontrado", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Instructor obtenido", instructor)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Crear instructor
// @Description Crea un nuevo instructor (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.InstructorRequest	true	"Datos del instructor"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [post]
func (c *InstructorController) Crear() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	var req models.InstructorRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	instructor, err := c.instructorService.Create(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al crear instructor", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Instructor creado correctamente", instructor)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Actualizar instructor
// @Description Actualiza un instructor (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID del instructor"
// @Param	body	body 	models.InstructorUpdateRequest	true	"Datos a actualizar"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [put]
func (c *InstructorController) Actualizar() {
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

	var req models.InstructorUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	instructor, err := c.instructorService.Update(id, req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al actualizar instructor", err)
		c.ServeJSON()
		return
	}

	if instructor == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Instructor no encontrado", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Instructor actualizado correctamente", instructor)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar instructor
// @Description Elimina un instructor (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID del instructor"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *InstructorController) Eliminar() {
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

	instructor, err := c.instructorService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar instructor", err)
		c.ServeJSON()
		return
	}
	if instructor == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Instructor no encontrado", nil)
		c.ServeJSON()
		return
	}

	err = c.instructorService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar instructor", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Instructor eliminado correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}
