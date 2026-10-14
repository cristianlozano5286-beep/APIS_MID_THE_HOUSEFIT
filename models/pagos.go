package models

import (
	"time"
)

type MetodoPago struct {
	ID          int       `json:"id" db:"id"`
	Nombre      string    `json:"nombre" db:"nombre"`
	Icono       string    `json:"icono" db:"icono"`
	Descripcion string    `json:"descripcion" db:"descripcion"`
	Disponible  bool      `json:"disponible" db:"disponible"`
}

type MetodoPagoRequest struct {
	Nombre      string `json:"nombre" binding:"required,min=2,max=100"`
	Icono       string `json:"icono" binding:"required,max=10"`
	Descripcion string `json:"descripcion"`
	Disponible  bool   `json:"disponible"`
}

type MetodoPagoUpdateRequest struct {
	Nombre      *string `json:"nombre" binding:"omitempty,min=2,max=100"`
	Icono       *string `json:"icono" binding:"omitempty,max=10"`
	Descripcion *string `json:"descripcion" binding:"omitempty"`
	Disponible  *bool   `json:"disponible" binding:"omitempty"`
}

type TransaccionRequest struct {
	UsuarioCorreo string  `json:"usuario_correo" binding:"required,email,max=150"`
	Concepto      string  `json:"concepto" binding:"required,min=2,max=200"`
	Monto         float64 `json:"monto" binding:"required,min=0.01"`
	MetodoPagoID  int     `json:"metodo_pago_id" binding:"required,min=1"`
}

type TransaccionResponse struct {
	Transaccion
	MetodoPago *MetodoPago `json:"metodo_pago"`
}