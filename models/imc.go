package models

import (
	"time"
)

type ResultadoIMC struct {
	ID           int       `json:"id" db:"id"`
	UsuarioCorreo *string  `json:"usuario_correo" db:"usuario_correo"`
	PesoKg       float64   `json:"peso_kg" db:"peso_kg"`
	EstaturaM    float64   `json:"estatura_m" db:"estatura_m"`
	IMC          float64   `json:"imc" db:"imc"`
	Categoria    string    `json:"categoria" db:"categoria"`
	Fecha        time.Time `json:"fecha" db:"fecha"`
}

type CalcularIMCRequest struct {
	UsuarioCorreo *string  `json:"usuario_correo" binding:"omitempty,email,max=150"`
	PesoKg        float64  `json:"peso_kg" binding:"required,min=1,max=500"`
	EstaturaM     float64  `json:"estatura_m" binding:"required,min=0.5,max=3"`
}

type CalcularIMCResponse struct {
	IMC       float64 `json:"imc"`
	Categoria string  `json:"categoria"`
	PesoIdeal struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"peso_ideal"`
	GuiaNutricional *GuiaNutricional `json:"guia_nutricional,omitempty"`
}





