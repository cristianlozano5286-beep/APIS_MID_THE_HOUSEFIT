package models

import (
	"time"
)

type GuiaNutricional struct {
	ID               int       `json:"id" db:"id"`
	CategoriaIMC     string    `json:"categoria_imc" db:"categoria_imc"`
	Titulo           string    `json:"titulo" db:"titulo"`
	Recomendaciones  string    `json:"recomendaciones" db:"recomendaciones"`
	EjemploComidas   string    `json:"ejemplo_comidas" db:"ejemplo_comidas"`
	CreadoEn         time.Time `json:"creado_en" db:"creado_en"`
}

type GuiaNutricionalRequest struct {
	CategoriaIMC    string `json:"categoria_imc" binding:"required,oneof='Bajo peso' 'Peso normal' Sobrepeso Obesidad"`
	Titulo          string `json:"titulo" binding:"required,min=2,max=200"`
	Recomendaciones string `json:"recomendaciones"`
	EjemploComidas  string `json:"ejemplo_comidas"`
}

type GuiaNutricionalUpdateRequest struct {
	CategoriaIMC    *string `json:"categoria_imc" binding:"omitempty,oneof='Bajo peso' 'Peso normal' Sobrepeso Obesidad"`
	Titulo          *string `json:"titulo" binding:"omitempty,min=2,max=200"`
	Recomendaciones *string `json:"recomendaciones" binding:"omitempty"`
	EjemploComidas  *string `json:"ejemplo_comidas" binding:"omitempty"`
}

const (
	CategoriaIMCBajoPeso   = "Bajo peso"
	CategoriaIMCPesoNormal = "Peso normal"
	CategoriaIMCSobrepeso  = "Sobrepeso"
	CategoriaIMCObesidad   = "Obesidad"
)