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