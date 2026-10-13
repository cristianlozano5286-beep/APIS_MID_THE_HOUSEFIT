package models

import (
	"time"
)

type Resena struct {
	ID           int       `json:"id" db:"id"`
	GimnasioID   int       `json:"gimnasio_id" db:"gimnasio_id"`
	UsuarioCorreo *string  `json:"usuario_correo" db:"usuario_correo"`
	UsuarioNombre string   `json:"usuario_nombre" db:"usuario_nombre"`
	Titulo       string    `json:"titulo" db:"titulo"`
	Cuerpo       string    `json:"cuerpo" db:"cuerpo"`
	Imagen       string    `json:"imagen" db:"imagen"`
	Puntuacion   int       `json:"puntuacion" db:"puntuacion"`
	Fecha        time.Time `json:"fecha" db:"fecha"`
}

type ResenaRequest struct {
	GimnasioID    int     `json:"gimnasio_id" binding:"required,min=1"`
	UsuarioCorreo *string `json:"usuario_correo" binding:"omitempty,email,max=150"`
	UsuarioNombre string  `json:"usuario_nombre" binding:"required,min=2,max=120"`
	Titulo        string  `json:"titulo" binding:"required,min=2,max=200"`
	Cuerpo        string  `json:"cuerpo" binding:"required"`
	Imagen        string  `json:"imagen" binding:"max=10"`
	Puntuacion    int     `json:"puntuacion" binding:"required,min=1,max=5"`
}

type ResenaUpdateRequest struct {
	Titulo     *string `json:"titulo" binding:"omitempty,min=2,max=200"`
	Cuerpo     *string `json:"cuerpo" binding:"omitempty"`
	Imagen     *string `json:"imagen" binding:"omitempty,max=10"`
	Puntuacion *int    `json:"puntuacion" binding:"omitempty,min=1,max=5"`
}

type ResenaResponse struct {
	Resena
	Gimnasio *Gimnasio `json:"gimnasio"`
}