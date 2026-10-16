package controllers

import (
	"strings"

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
