package models

import (
	"time"
)

type Producto struct {
	ID              int       `json:"id" db:"id"`
	Nombre          string    `json:"nombre" db:"nombre"`
	Categoria       string    `json:"categoria" db:"categoria"`
	Precio          float64   `json:"precio" db:"precio"`
	Stock           int       `json:"stock" db:"stock"`
	Descripcion     string    `json:"descripcion" db:"descripcion"`
	GuiaNutricional string    `json:"guia_nutricional" db:"guia_nutricional"`
	Imagen          string    `json:"imagen" db:"imagen"`
	CreadoEn        time.Time `json:"creado_en" db:"creado_en"`
}

type ProductoRequest struct {
	Nombre           string  `json:"nombre" binding:"required,min=2,max=150"`
	Categoria        string  `json:"categoria" binding:"required,oneof=Suplemento Implemento 'Ropa deportiva'"`
	Precio           float64 `json:"precio" binding:"required,min=0"`
	Stock            int     `json:"stock" binding:"required,min=0"`
	Descripcion      string  `json:"descripcion"`
	GuiaNutricional  string  `json:"guia_nutricional"`
	Imagen           string  `json:"imagen" binding:"max=10"`
}

type ProductoUpdateRequest struct {
	Nombre           *string  `json:"nombre" binding:"omitempty,min=2,max=150"`
	Categoria        *string  `json:"categoria" binding:"omitempty,oneof=Suplemento Implemento 'Ropa deportiva'"`
	Precio           *float64 `json:"precio" binding:"omitempty,min=0"`
	Stock            *int     `json:"stock" binding:"omitempty,min=0"`
	Descripcion      *string  `json:"descripcion" binding:"omitempty"`
	GuiaNutricional  *string  `json:"guia_nutricional" binding:"omitempty"`
	Imagen           *string  `json:"imagen" binding:"omitempty,max=10"`
}

const (
	CategoriaSuplemento     = "Suplemento"
	CategoriaImplemento     = "Implemento"
	CategoriaRopaDeportiva  = "Ropa deportiva"
)