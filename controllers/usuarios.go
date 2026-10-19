package controllers

import (
	beego "github.com/beego/beego/v2/server/web"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"
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

// GET /api/v1/resumen
// los numeros que muestra el dashboard del administrador
func (c *UsuariosController) Resumen() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	colecciones := []string{"usuarios", "gimnasios", "instructores", "productos", "clases", "reservas-clases", "resenas", "noticias"}

	resumen := map[string]interface{}{}

	for _, coleccion := range colecciones {
		var filas []map[string]interface{}
		if !listarColeccion(&c.Controller, coleccion, &filas) {
			return
		}
		resumen[coleccion] = len(filas)
	}

	var transacciones []models.Transaccion
	if !listarColeccion(&c.Controller, "transacciones", &transacciones) {
		return
	}

	ingresos := 0.0
	for _, t := range transacciones {
		if t.Estado == "Aprobado" {
			ingresos += t.Monto
		}
	}
	resumen["transacciones"] = len(transacciones)
	resumen["ingresos"] = ingresos

	responder(&c.Controller, 200, resumen)
}
