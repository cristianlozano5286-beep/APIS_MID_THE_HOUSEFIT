package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type ClaseController struct {
	beego.Controller
	claseService   *services.ClaseService
	reservaService *services.ReservaService
}

func (c *ClaseController) Prepare() {
	c.claseService = services.NewClaseService()
	c.reservaService = services.NewReservaService()
}

// @Title Listar clases
// @Description Lista clases disponibles con filtros y paginación
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Param	gimnasio_id	query	int	false	"Filtrar por gimnasio"
// @Param	tipo_clase	query	string	false	"Filtrar por tipo de clase"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 500 {object} models.RespuestaError
// @router / [get]
func (c *ClaseController) Listar() {
	pagina, _ := strconv.Atoi(c.GetString("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.GetString("por_pagina", "10"))
	gimnasioIDStr := c.GetString("gimnasio_id")
	tipoClase := c.GetString("tipo_clase")

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
	if tipoClase != "" {
		filters["tipo_clase"] = tipoClase
	}

	clases, total, err := c.claseService.List(filters, pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar clases", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Clases obtenidas", clases, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Ver clase
// @Description Obtiene una clase por ID
// @Param	id	path	int	true	"ID de la clase"
// @Success 200 {object} models.RespuestaAPI
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [get]
func (c *ClaseController) Ver() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("ID inválido", err)
		c.ServeJSON()
		return
	}

	clase, err := c.claseService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener clase", err)
		c.ServeJSON()
		return
	}

	if clase == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Clase no encontrada", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Clase obtenida", clase)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Crear clase
// @Description Crea una nueva clase (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.ClaseDisponibleRequest	true	"Datos de la clase"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router / [post]
func (c *ClaseController) Crear() {
	rol := c.Ctx.Input.GetData("usuario_rol")
	if rol != models.RolAdministrador {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		c.Data["json"] = models.RespuestaErrorGeneral("Acceso denegado", nil)
		c.ServeJSON()
		return
	}

	var req models.ClaseDisponibleRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	clase, err := c.claseService.Create(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al crear clase", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Clase creada correctamente", clase)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar clase
// @Description Elimina una clase (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la clase"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *ClaseController) Eliminar() {
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

	clase, err := c.claseService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar clase", err)
		c.ServeJSON()
		return
	}
	if clase == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Clase no encontrada", nil)
		c.ServeJSON()
		return
	}

	err = c.claseService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar clase", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Clase eliminada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Reservar clase
// @Description Reserva una plaza en una clase
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.ReservaClaseRequest	true	"Datos de la reserva"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 409 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /reservar [post]
func (c *ClaseController) Reservar() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
	if usuarioCorreo == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
		c.ServeJSON()
		return
	}

	var req models.ReservaClaseRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	req.UsuarioCorreo = usuarioCorreo.(string)

	reserva, err := c.reservaService.Create(req)
	if err != nil {
		if err.Error() == "no hay cupos disponibles para esta fecha" {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
			c.Data["json"] = models.RespuestaErrorGeneral("No hay cupos disponibles", err)
			c.ServeJSON()
			return
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al reservar clase", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Clase reservada correctamente", reserva)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Mis reservas
// @Description Obtiene las reservas del usuario autenticado
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	pagina	query	int	false	"Página (default: 1)"
// @Param	por_pagina	query	int	false	"Por página (default: 10)"
// @Success 200 {object} models.RespuestaPaginada
// @Failure 401 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /mis-reservas [get]
func (c *ClaseController) MisReservas() {
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

	reservas, total, err := c.reservaService.ListByUsuario(usuarioCorreo.(string), pagina, porPagina)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al obtener reservas", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaPaginadaExitosa("Reservas obtenidas", reservas, pagina, porPagina, total)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Cancelar reserva
// @Description Cancela una reserva del usuario
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID de la reserva"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /reservas/:id [delete]
func (c *ClaseController) CancelarReserva() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
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

	err = c.reservaService.Cancelar(id, usuarioCorreo.(string))
	if err != nil {
		if err.Error() == "reserva no encontrada o no pertenece al usuario" {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
			c.Data["json"] = models.RespuestaErrorGeneral("Reserva no encontrada", err)
			c.ServeJSON()
			return
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al cancelar reserva", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Reserva cancelada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}