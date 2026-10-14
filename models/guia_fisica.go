package models

import (
	"time"
)

type GuiaFisica struct {
	ID               int       `json:"id" db:"id"`
	Titulo           string    `json:"titulo" db:"titulo"`
	Metodo           string    `json:"metodo" db:"metodo"`
	Contenido        string    `json:"contenido" db:"contenido"`
	FechaPublicacion time.Time `json:"fecha_publicacion" db:"fecha_publicacion"`
	CreadoEn         time.Time `json:"creado_en" db:"creado_en"`
}

type GuiaFisicaRequest struct {
	Titulo           string `json:"titulo" binding:"required,min=2,max=200"`
	Metodo           string `json:"metodo" binding:"required,max=100"`
	Contenido        string `json:"contenido"`
	FechaPublicacion string `json:"fecha_publicacion"` // YYYY-MM-DD
}

type GuiaFisicaUpdateRequest struct {
	Titulo           *string `json:"titulo" binding:"omitempty,min=2,max=200"`
	Metodo           *string `json:"metodo" binding:"omitempty,max=100"`
	Contenido        *string `json:"contenido" binding:"omitempty"`
	FechaPublicacion *string `json:"fecha_publicacion" binding:"omitempty"`
}