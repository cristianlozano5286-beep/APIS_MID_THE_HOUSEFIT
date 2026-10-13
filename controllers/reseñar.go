package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type ResenaController struct {
	beego.Controller
	resenaService *services.ResenaService
}

func (c *ResenaController) Prepare() {
	c.resenaService = services.NewResenaService()
}

// @Title Listar reseñas
// @Description Lista reseñas de un gimnasio
// @Param	gimnasio_id	query	int	true	"ID del gimnasio"
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 400 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *ResenaController) Listar() {
	gimnasioIDStr := c.GetString("gimnasio_id")
	if gimnasioIDStr == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("gimnasio_id requerido", nil)
		c.ServeJSON()
		return
	}

	gimnasioID, err := strconv.Atoi(gimnasioIDStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("gimnasio_id inválido", err)
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

	resenas, total, err := c.resenaService.ListByGimnasio(gimnasioID, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar reseñas", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Reseñas obtenidas", resenas, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Agregar reseña
// @Description Agrega una reseña a un gimnasio
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.ResenaRequest	true	"Datos de la reseña"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [post]
func (c *ResenaController) Agregar() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
	if usuarioCorreo == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
		c.ServeJSON()
		return
	}

	var req models.ResenaRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	authService := services.NewAuthService()
	usuario, _ := authService.ObtenerUsuarioPorCorreo(usuarioCorreo.(string))
	if usuario != nil {
		req.UsuarioCorreo = &usuario.Correo
		req.UsuarioNombre = usuario.Nombre
	} else {
		req.UsuarioNombre = "Usuario"
	}

	resena, err := c.resenaService.Create(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al agregar reseña", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Reseña agregada correctamente", resena)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar reseña
// @Description Elimina una reseña (solo autor o administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la reseña"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *ResenaController) Eliminar() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
	usuarioRol := c.Ctx.Input.GetData("usuario_rol")
	if usuarioCorreo == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
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

	resena, err := c.resenaService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar reseña", err)
		c.ServeJSON()
		return
	}
	if resena == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Reseña no encontrada", nil)
		c.ServeJSON()
		return
	}

	if resena.Resena.UsuarioCorreo != nil && *resena.Resena.UsuarioCorreo != usuarioCorreo.(string) && usuarioRol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("No autorizado para eliminar esta reseña", nil)
		c.ServeJSON()
		return
	}

	err = c.resenaService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar reseña", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Reseña eliminada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Obtener promedio
// @Description Obtiene el promedio de calificaciones de un gimnasio
// @Param	gimnasio_id	query	int	true	"ID del gimnasio"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /promedio [get]
func (c *ResenaController) Promedio() {
	gimnasioIDStr := c.GetString("gimnasio_id")
	if gimnasioIDStr == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("gimnasio_id requerido", nil)
		c.ServeJSON()
		return
	}

	gimnasioID, err := strconv.Atoi(gimnasioIDStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("gimnasio_id inválido", err)
		c.ServeJSON()
		return
	}

	promedio, total, err := c.resenaService.GetPromedioGimnasio(gimnasioID)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener promedio", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Promedio obtenido", map[string]interface{}{
		"promedio":  promedio,
		"total":     total,
	})
	c.Data["json"] = response
	c.ServeJSON()
}