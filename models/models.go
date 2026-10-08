package models

import "time"


type ClaseDisponible struct {
	ID      int     `json:"id"`
	GimnasioID    int     `json:"gimnasio_id"`
	TipoClase  string  `json:"tipo_clase"`
	InstructorNombre string  `json:"instructor_nombre"`
	InstructorID     int     `json:"instructor_id,omitempty"`
	Hora   string  `json:"hora"`
	Lugar string `json:"lugar"`
	Cupos    int     `json:"cupos"`
	Duracion    string  `json:"duracion"`
	Precio   float64 `json:"precio"`
}

type ReservaClase struct {
	ID            int       `json:"id"`
	UsuarioCorreo string    `json:"usuario_correo"`
	ClaseID       int       `json:"clase_id"`
	Fecha         string    `json:"fecha"`
	Personas      int       `json:"personas"`
	Estado        string    `json:"estado"` 
	CreadoEn      time.Time `json:"creado_en"`
}
