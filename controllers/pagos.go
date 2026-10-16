package controllers

import (
<<<<<<< HEAD
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
=======
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type PagoController struct {
	beego.Controller
	pagoService *services.PagoService
}

func (c *PagoController) Prepare() {
	c.pagoService = services.NewPagoService()
}

// @Title Listar métodos de pago
// @Description Lista los métodos de pago disponibles
// @Success 200 {object} models.RespuestaAPI
// @Failure 500 {object} models.RespuestaError
// @router /metodos [get]
func (c *PagoController) ListarMetodos() {
	metodos, err := c.pagoService.ListarMetodos()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al listar métodos de pago", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Métodos de pago obtenidos", metodos)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Procesar pago
// @Description Procesa un pago (simulado)
// @Param	Authorization	header	string	true	"Bearer token"
// @Param	body	body 	models.ProcesarPagoRequest	true	"Datos del pago"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @Failure 500 {object} models.RespuestaError
// @router /procesar [post]
func (c *PagoController) ProcesarPago() {
	usuarioCorreo := c.Ctx.Input.GetData("usuario_correo")
	if usuarioCorreo == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
		c.ServeJSON()
		return
	}

	var req models.ProcesarPagoRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	req.UsuarioCorreo = usuarioCorreo.(string)

	resultado, err := c.pagoService.ProcesarPago(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al procesar pago", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Pago procesado", resultado)
	c.Data["json"] = response
	c.ServeJSON()
}
>>>>>>> c329136a84f67cabe8f6c766b9b85bb8c59ea932
