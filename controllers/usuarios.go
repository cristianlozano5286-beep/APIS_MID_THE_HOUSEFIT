package controllers

import (
	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
	"thehousefit_mid/services"
)

type UsuariosController struct {
	beego.Controller
}

// GET /api/v1/usuarios
func (c *UsuariosController) Listar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	var usuarios []models.Usuario
	if !listarColeccion(&c.Controller, "usuarios", &usuarios) {
		return
	}

	// la contraseña nunca sale del MID
	for i := range usuarios {
		usuarios[i].Password = ""
	}

	responder(&c.Controller, 200, usuarios)
}

// PUT /api/v1/usuarios/:id/rol
func (c *UsuariosController) CambiarRol() {
	sesion, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	id, valido := leerId(&c.Controller)
	if !valido {
		return
	}

	var cambio models.CambioRol

	if !leerCuerpo(&c.Controller, &cambio) {
		return
	}

	if cambio.Rol != "Administrador" && cambio.Rol != "Entrenador" && cambio.Rol != "Usuario" {
		responderError(&c.Controller, 400, "El rol debe ser Administrador, Entrenador o Usuario.")
		return
	}

	// para no quedarse sin administradores
	if id == sesion.ID {
		responderError(&c.Controller, 400, "No puedes cambiar tu propio rol.")
		return
	}

	cambios := map[string]interface{}{"rol": cambio.Rol}

	var actualizado models.Usuario
	err := services.NuevoServicioCrud().Actualizar("usuarios", id, cambios, &actualizado)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	actualizado.Password = ""
	responder(&c.Controller, 200, actualizado)
}
