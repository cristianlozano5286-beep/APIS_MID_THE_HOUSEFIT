package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
	"thehousefit_mid/services"
)

type PagosController struct {
	beego.Controller
}

// GET /api/v1/pagos/metodos
func (c *PagosController) ListarMetodos() {
	var metodos []models.MetodoPago
	if !listarColeccion(&c.Controller, "metodos-pago", &metodos) {
		return
	}

	// solo los que se pueden usar
	disponibles := []models.MetodoPago{}
	for _, m := range metodos {
		if m.Disponible {
			disponibles = append(disponibles, m)
		}
	}

	responder(&c.Controller, 200, disponibles)
}

// POST /api/v1/pagos
// el pago es simulado: no hay una pasarela de pago de verdad
func (c *PagosController) Pagar() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	var solicitud models.SolicitudPago

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	if strings.TrimSpace(solicitud.Concepto) == "" {
		responderError(&c.Controller, 400, "El concepto es obligatorio.")
		return
	}

	if solicitud.Monto <= 0 {
		responderError(&c.Controller, 400, "El monto debe ser mayor a cero.")
		return
	}

	// el metodo de pago tiene que existir y estar disponible
	var metodos []models.MetodoPago
	if !listarColeccion(&c.Controller, "metodos-pago", &metodos) {
		return
	}

	metodoBueno := false
	for _, m := range metodos {
		if m.Id == solicitud.MetodoPagoId && m.Disponible {
			metodoBueno = true
		}
	}
	if !metodoBueno {
		responderError(&c.Controller, 400, "Ese método de pago no está disponible.")
		return
	}

	nueva := models.Transaccion{
		Usuario:      sesion.Nombre,
		Concepto:     strings.TrimSpace(solicitud.Concepto),
		Monto:        solicitud.Monto,
		MetodoPagoId: solicitud.MetodoPagoId,
		Estado:       "Aprobado",
		Fecha:        hoy(),
	}

	var creada map[string]interface{}
	err := services.NuevoServicioCrud().Crear("transacciones", nueva, &creada)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 201, creada)
}

// GET /api/v1/pagos/transacciones
// los pagos de la persona que inicio sesion
func (c *PagosController) MisTransacciones() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	var transacciones []models.Transaccion
	if !listarColeccion(&c.Controller, "transacciones", &transacciones) {
		return
	}

	mias := []models.Transaccion{}
	for _, t := range transacciones {
		if t.Usuario == sesion.Nombre {
			mias = append(mias, t)
		}
	}

	responder(&c.Controller, 200, mias)
}
