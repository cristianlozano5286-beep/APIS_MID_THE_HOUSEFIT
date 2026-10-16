package models

import "time"

type Noticia struct {
	ID        int       `json:"id" db:"id"`
	Titulo    string    `json:"titulo" db:"titulo"`
	Slug      string    `json:"slug" db:"slug"`
	Resumen   string    `json:"resumen" db:"resumen"`
	Contenido string    `json:"contenido" db:"contenido"`
	Imagen    string    `json:"imagen" db:"imagen"`
	Tipo      string    `json:"tipo" db:"tipo"`
	Activa    bool      `json:"activa" db:"activa"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
}

type NoticiaRequest struct {
	Titulo    string `json:"titulo" binding:"required,min=2,max=200"`
	Slug      string `json:"slug" binding:"required,min=2,max=200"`
	Resumen   string `json:"resumen" binding:"required,min=2,max=500"`
	Contenido string `json:"contenido" binding:"required"`
	Imagen    string `json:"imagen" binding:"max=200"`
	Tipo      string `json:"tipo" binding:"required,max=50"`
	Activa    bool   `json:"activa"`
}

type NoticiaUpdateRequest struct {
	Titulo    *string `json:"titulo" binding:"omitempty,min=2,max=200"`
	Slug      *string `json:"slug" binding:"omitempty,min=2,max=200"`
	Resumen   *string `json:"resumen" binding:"omitempty,min=2,max=500"`
	Contenido *string `json:"contenido" binding:"omitempty"`
	Imagen    *string `json:"imagen" binding:"omitempty,max=200"`
	Tipo      *string `json:"tipo" binding:"omitempty,max=50"`
	Activa    *bool   `json:"activa"`
}
