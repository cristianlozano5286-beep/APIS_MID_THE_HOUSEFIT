package models

import (
	"time"
)

type Instructor struct {
	ID           int       `json:"id" db:"id"`
	Nombre       string    `json:"nombre" db:"nombre"`
	Correo       string    `json:"correo" db:"correo"`
	Telefono     string    `json:"telefono" db:"telefono"`
	Especialidad string    `json:"especialidad" db:"especialidad"`
	Experiencia  int       `json:"experiencia" db:"experiencia"`
	Biografia    string    `json:"biografia" db:"biografia"`
	Imagen       string    `json:"imagen" db:"imagen"`
	GimnasioID   *int      `json:"gimnasio_id" db:"gimnasio_id"`
	Activo       bool      `json:"activo" db:"activo"`
	CreadoEn     time.Time `json:"creado_en" db:"creado_en"`
}

type InstructorRequest struct {
	Nombre       string `json:"nombre" binding:"required,min=2,max=120"`
	Correo       string `json:"correo" binding:"required,email,max=150"`
	Telefono     string `json:"telefono" binding:"required,max=30"`
	Especialidad string `json:"especialidad" binding:"required,max=120"`
	Experiencia  int    `json:"experiencia" binding:"required,min=0"`
	Biografia    string `json:"biografia"`
	Imagen       string `json:"imagen" binding:"max=10"`
	GimnasioID   *int   `json:"gimnasio_id"`
	Activo       bool   `json:"activo"`
}

type InstructorUpdateRequest struct {
	Nombre       *string `json:"nombre" binding:"omitempty,min=2,max=120"`
	Correo       *string `json:"correo" binding:"omitempty,email,max=150"`
	Telefono     *string `json:"telefono" binding:"omitempty,max=30"`
	Especialidad *string `json:"especialidad" binding:"omitempty,max=120"`
	Experiencia  *int    `json:"experiencia" binding:"omitempty,min=0"`
	Biografia    *string `json:"biografia" binding:"omitempty"`
	Imagen       *string `json:"imagen" binding:"omitempty,max=10"`
	GimnasioID   *int    `json:"gimnasio_id"`
	Activo       *bool   `json:"activo"`
}