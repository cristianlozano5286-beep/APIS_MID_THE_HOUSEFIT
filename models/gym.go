package models

import (
	"time"
)

type Gimnasio struct {
	ID          int       `json:"id" db:"id"`
	Nombre      string    `json:"nombre" db:"nombre"`
	Direccion   string    `json:"direccion" db:"direccion"`
	Barrio      string    `json:"barrio" db:"barrio"`
	Ciudad      string    `json:"ciudad" db:"ciudad"`
	Pais        string    `json:"pais" db:"pais"`
	Telefono    string    `json:"telefono" db:"telefono"`
	Horario     string    `json:"horario" db:"horario"`
	Descripcion string    `json:"descripcion" db:"descripcion"`
	Imagen      string    `json:"imagen" db:"imagen"`
	Precio      float64   `json:"precio" db:"precio"`
	Moneda      string    `json:"moneda" db:"moneda"`
	Activo      bool      `json:"activo" db:"activo"`
	Lat         float64   `json:"lat" db:"lat"`
	Lng         float64   `json:"lng" db:"lng"`
	CreadoEn    time.Time `json:"creado_en" db:"creado_en"`
}

type GimnasioFoto struct {
	ID        int       `json:"id" db:"id"`
	GimnasioID int      `json:"gimnasio_id" db:"gimnasio_id"`
	URL       string    `json:"url" db:"url"`
	Orden     int       `json:"orden" db:"orden"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
}

type GimnasioServicio struct {
	ID        int       `json:"id" db:"id"`
	GimnasioID int      `json:"gimnasio_id" db:"gimnasio_id"`
	Servicio  string    `json:"servicio" db:"servicio"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
}

type GimnasioCompleto struct {
	Gimnasio
	Fotos     []GimnasioFoto     `json:"fotos"`
	Servicios []GimnasioServicio `json:"servicios"`
}

type GimnasioRequest struct {
	Nombre      string   `json:"nombre" binding:"required,min=2,max=120"`
	Direccion   string   `json:"direccion" binding:"required,max=200"`
	Barrio      string   `json:"barrio" binding:"required,max=100"`
	Ciudad      string   `json:"ciudad" binding:"required,max=100"`
	Pais        string   `json:"pais" binding:"max=100"`
	Telefono    string   `json:"telefono" binding:"required,max=30"`
	Horario     string   `json:"horario" binding:"required,max=150"`
	Descripcion string   `json:"descripcion" binding:"required"`
	Imagen      string   `json:"imagen" binding:"max=10"`
	Precio      float64  `json:"precio" binding:"required,min=0"`
	Moneda      string   `json:"moneda" binding:"max=10"`
	Activo      bool     `json:"activo"`
	Lat         float64  `json:"lat"`
	Lng         float64  `json:"lng"`
	Fotos       []string `json:"fotos"`
	Servicios   []string `json:"servicios"`
}

type GimnasioUpdateRequest struct {
	Nombre      *string  `json:"nombre" binding:"omitempty,min=2,max=120"`
	Direccion   *string  `json:"direccion" binding:"omitempty,max=200"`
	Barrio      *string  `json:"barrio" binding:"omitempty,max=100"`
	Ciudad      *string  `json:"ciudad" binding:"omitempty,max=100"`
	Pais        *string  `json:"pais" binding:"omitempty,max=100"`
	Telefono    *string  `json:"telefono" binding:"omitempty,max=30"`
	Horario     *string  `json:"horario" binding:"omitempty,max=150"`
	Descripcion *string  `json:"descripcion" binding:"omitempty"`
	Imagen      *string  `json:"imagen" binding:"omitempty,max=10"`
	Precio      *float64 `json:"precio" binding:"omitempty,min=0"`
	Moneda      *string  `json:"moneda" binding:"omitempty,max=10"`
	Activo      *bool    `json:"activo"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	Fotos       []string `json:"fotos"`
	Servicios   []string `json:"servicios"`
}