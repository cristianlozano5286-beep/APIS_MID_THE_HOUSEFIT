package controllers

import (
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"
)

type ClasesController struct {
	beego.Controller
}

// GET /api/v1/clases?gimnasioId=1&tipoClase=Yoga
func (c *ClasesController) Listar() {
	var clases []models.Clase
	if !listarColeccion(&c.Controller, "clases", &clases) {
		return
	}

	gimnasioId := c.GetString("gimnasioId")
	tipoClase := strings.TrimSpace(c.GetString("tipoClase"))

	resultado := []models.Clase{}
	for _, clase := range clases {
		if gimnasioId != "" && strconv.Itoa(clase.GimnasioId) != gimnasioId {
			continue
		}
		if tipoClase != "" && clase.TipoClase != tipoClase {
			continue
		}
		resultado = append(resultado, clase)
	}

	responder(&c.Controller, 200, resultado)
}

// POST /api/v1/clases
func (c *ClasesController) Crear() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	var clase models.Clase

	if !leerCuerpo(&c.Controller, &clase) {
		return
	}

	if clase.GimnasioId == 0 || strings.TrimSpace(clase.Instructor) == "" || strings.TrimSpace(clase.Hora) == "" {
		responderError(&c.Controller, 400, "El gimnasio, el instructor y la hora son obligatorios.")
		return
	}

	crearEnColeccion(&c.Controller, "clases", clase)
}

// DELETE /api/v1/clases/:id
func (c *ClasesController) Borrar() {
	_, valido := exigirRol(&c.Controller, "Administrador", "Entrenador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "clases")
}

// POST /api/v1/clases/:id/reservar
func (c *ClasesController) Reservar() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	id, valido := leerId(&c.Controller)
	if !valido {
		return
	}
	claseId, _ := strconv.Atoi(id)

	var solicitud models.SolicitudReserva

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	if solicitud.Personas < 1 {
		responderError(&c.Controller, 400, "Debe reservar al menos una persona.")
		return
	}

	fecha, err := time.Parse("2006-01-02", solicitud.Fecha)
	if err != nil {
		responderError(&c.Controller, 400, "La fecha debe tener el formato AAAA-MM-DD.")
		return
	}
	if fecha.Format("2006-01-02") < hoy() {
		responderError(&c.Controller, 400, "No se puede reservar en una fecha pasada.")
		return
	}

	crud := services.NuevoServicioCrud()

	var clase models.Clase
	err = crud.Obtener("clases", id, &clase)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	var reservas []models.Reserva
	if !listarColeccion(&c.Controller, "reservas-clases", &reservas) {
		return
	}

	// cuantos cupos ya estan ocupados ese dia
	ocupados := 0
	for _, r := range reservas {
		if r.ClaseId == claseId && r.Fecha == solicitud.Fecha && r.Estado == "Confirmada" {
			ocupados += r.Personas
		}
	}

	if ocupados+solicitud.Personas > clase.Cupos {
		responderError(&c.Controller, 409, "No hay cupos suficientes para ese día.")
		return
	}

	nueva := models.Reserva{
		Usuario:  sesion.Nombre,
		ClaseId:  claseId,
		Fecha:    solicitud.Fecha,
		Personas: solicitud.Personas,
		Estado:   "Confirmada",
	}

	crearEnColeccion(&c.Controller, "reservas-clases", nueva)
}

// GET /api/v1/reservas
// las reservas de la persona que inicio sesion
func (c *ClasesController) MisReservas() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	var reservas []models.Reserva
	if !listarColeccion(&c.Controller, "reservas-clases", &reservas) {
		return
	}

	mias := []models.Reserva{}
	for _, r := range reservas {
		if r.Usuario == sesion.Nombre {
			mias = append(mias, r)
		}
	}

	responder(&c.Controller, 200, mias)
}

// PUT /api/v1/reservas/:id/cancelar
func (c *ClasesController) CancelarReserva() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	id, valido := leerId(&c.Controller)
	if !valido {
		return
	}

	crud := services.NuevoServicioCrud()

	var reserva models.Reserva
	err := crud.Obtener("reservas-clases", id, &reserva)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if reserva.Usuario != sesion.Nombre {
		responderError(&c.Controller, 403, "Esta reserva no es tuya.")
		return
	}

	if reserva.Estado == "Cancelada" {
		responderError(&c.Controller, 400, "Esta reserva ya estaba cancelada.")
		return
	}

	cambios := map[string]interface{}{"estado": "Cancelada"}

	var actualizada models.Reserva
	err = crud.Actualizar("reservas-clases", id, cambios, &actualizada)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, actualizada)
}
