package models

import (
	"time"
)

type CodigoVerificacion struct {
	ID        int       `json:"id" db:"id"`
	Correo    string    `json:"correo" db:"correo"`
	Codigo    string    `json:"codigo" db:"codigo"`
	Tipo      string    `json:"tipo" db:"tipo"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
	ExpiraEn  time.Time `json:"expira_en" db:"expira_en"`
	Usado     bool      `json:"usado" db:"usado"`
}

type CodigoVerificacionRequest struct {
	Correo string `json:"correo" binding:"required,email,max=150"`
	Tipo   string `json:"tipo" binding:"required,oneof=verificacion recuperacion"`
}

type VerificarCodigoRequest struct {
	Correo string `json:"correo" binding:"required,email,max=150"`
	Codigo string `json:"codigo" binding:"required,min=6,max=6"`
	Tipo   string `json:"tipo" binding:"required,oneof=verificacion recuperacion"`
}

type RecuperarPasswordRequest struct {
	Correo string `json:"correo" binding:"required,email,max=150"`
}

type CambiarPasswordRequest struct {
	Correo            string `json:"correo" binding:"required,email,max=150"`
	Codigo            string `json:"codigo" binding:"required,min=6,max=6"`
	NuevaPassword     string `json:"nueva_password" binding:"required,min=6,max=255"`
	ConfirmarPassword string `json:"confirmar_password" binding:"required"`
}

type Sesion struct {
	Token     string    `json:"token" db:"token"`
	Correo    string    `json:"correo" db:"correo"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
	ExpiraEn  time.Time `json:"expira_en" db:"expira_en"`
}

const (
	TipoCodigoVerificacion = "verificacion"
	TipoCodigoRecuperacion = "recuperacion"
	LongitudCodigo         = 6
	ExpiracionMinutos      = 15
)