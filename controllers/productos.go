package controllers

import (
<<<<<<< HEAD
	"strings"
=======
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"
>>>>>>> c329136a84f67cabe8f6c766b9b85bb8c59ea932

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
)

type ProductosController struct {
	beego.Controller
}

// GET /api/v1/productos
func (c *ProductosController) Listar() {
	var productos []models.Producto
	if !listarColeccion(&c.Controller, "productos", &productos) {
		return
	}

	responder(&c.Controller, 200, productos)
}

// POST /api/v1/productos
func (c *ProductosController) Crear() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	var producto models.Producto

	if !leerCuerpo(&c.Controller, &producto) {
		return
	}

	if strings.TrimSpace(producto.Nombre) == "" {
		responderError(&c.Controller, 400, "El nombre es obligatorio.")
		return
	}

	crearEnColeccion(&c.Controller, "productos", producto)
}

// PUT /api/v1/productos/:id
func (c *ProductosController) Actualizar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	actualizarEnColeccion(&c.Controller, "productos")
}

// DELETE /api/v1/productos/:id
func (c *ProductosController) Borrar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "productos")
}
<<<<<<< HEAD
=======

// @Title Actualizar producto
// @Description Actualiza un producto (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID del producto"
// @Param	body	body 	models.ProductoUpdateRequest	true	"Datos a actualizar"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [put]
func (c *ProductoController) Actualizar() {
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

	var req models.ProductoUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	producto, err := c.productoService.Update(id, req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al actualizar producto", err)
		c.ServeJSON()
		return
	}

	if producto == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Producto no encontrado", nil)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Producto actualizado correctamente", producto)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Eliminar producto
// @Description Elimina un producto (solo administrador)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	id	path	int	true	"ID del producto"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @Failure 403 {object} models.RespuestaError
// @Failure 404 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /:id [delete]
func (c *ProductoController) Eliminar() {
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

	producto, err := c.productoService.GetByID(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar producto", err)
		c.ServeJSON()
		return
	}
	if producto == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Producto no encontrado", nil)
		c.ServeJSON()
		return
	}

	err = c.productoService.Delete(id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al eliminar producto", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Producto eliminado correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}
>>>>>>> c329136a84f67cabe8f6c766b9b85bb8c59ea932
