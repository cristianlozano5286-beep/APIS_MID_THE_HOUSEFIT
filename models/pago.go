package models

type MetodoPago struct {
	Id          int    `json:"id,omitempty"`
	Nombre      string `json:"nombre"`
	Icono       string `json:"icono"`
	Descripcion string `json:"descripcion"`
	Disponible  bool   `json:"disponible"`
}

type Transaccion struct {
	Id           int     `json:"id,omitempty"`
	Usuario      string  `json:"usuario"`
	Concepto     string  `json:"concepto"`
	Monto        float64 `json:"monto"`
	MetodoPagoId int     `json:"metodoPagoId"`
	Estado       string  `json:"estado"`
	Fecha        string  `json:"fecha,omitempty"`
}

// lo que manda el cliente para pagar
type SolicitudPago struct {
	Concepto     string  `json:"concepto"`
	Monto        float64 `json:"monto"`
	MetodoPagoId int     `json:"metodoPagoId"`
}
