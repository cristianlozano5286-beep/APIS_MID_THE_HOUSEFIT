package models

import (
	"time"
)

type ClaseDisponible struct {
	ID               int       `json:"id" db:"id"`
	GimnasioID       int       `json:"gimnasio_id" db:"gimnasio_id"`
	TipoClase        string    `json:"tipo_clase" db:"tipo_clase"`
	InstructorNombre string    `json:"instructor_nombre" db:"instructor_nombre"`
	InstructorID     *int      `json:"instructor_id" db:"instructor_id"`
	Hora             string    `json:"hora" db:"hora"`
	Cupos            int       `json:"cupos" db:"cupos"`
	Duracion         string    `json:"duracion" db:"duracion"`
	Precio           float64   `json:"precio" db:"precio"`
	CreadoEn         time.Time `json:"creado_en" db:"creado_en"`
}

type ClaseDisponibleRequest struct {
	GimnasioID       int     `json:"gimnasio_id" binding:"required,min=1"`
	TipoClase        string  `json:"tipo_clase" binding:"required,oneof=Spinning CrossFit Yoga Funcional Boxeo Pilates"`
	InstructorNombre string  `json:"instructor_nombre" binding:"required,min=2,max=120"`
	InstructorID     *int    `json:"instructor_id" binding:"omitempty,min=1"`
	Hora             string  `json:"hora" binding:"required,max=20"`
	Cupos            int     `json:"cupos" binding:"required,min=0"`
	Duracion         string  `json:"duracion" binding:"required,max=20"`
	Precio           float64 `json:"precio" binding:"required,min=0"`
}

type ClaseDisponibleUpdateRequest struct {
	GimnasioID       *int     `json:"gimnasio_id" binding:"omitempty,min=1"`
	TipoClase        *string  `json:"tipo_clase" binding:"omitempty,oneof=Spinning CrossFit Yoga Funcional Boxeo Pilates"`
	InstructorNombre *string  `json:"instructor_nombre" binding:"omitempty,min=2,max=120"`
	InstructorID     *int     `json:"instructor_id" binding:"omitempty,min=1"`
	Hora             *string  `json:"hora" binding:"omitempty,max=20"`
	Cupos            *int     `json:"cupos" binding:"omitempty,min=0"`
	Duracion         *string  `json:"duracion" binding:"omitempty,max=20"`
	Precio           *float64 `json:"precio" binding:"omitempty,min=0"`
}

const (
	TipoClaseSpinning   = "Spinning"
	TipoClaseCrossFit   = "CrossFit"
	TipoClaseYoga       = "Yoga"
	TipoClaseFuncional  = "Funcional"
	TipoClaseBoxeo      = "Boxeo"
	TipoClasePilates    = "Pilates"
)