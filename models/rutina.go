package models

import (
	"time"
)

type Rutina struct {
	ID                 int       `json:"id" db:"id"`
	TipoEntrenamiento  string    `json:"tipo_entrenamiento" db:"tipo_entrenamiento"`
	Nombre             string    `json:"nombre" db:"nombre"`
	Nivel              string    `json:"nivel" db:"nivel"`
	Descripcion        string    `json:"descripcion" db:"descripcion"`
	DuracionSemanas    int       `json:"duracion_semanas" db:"duracion_semanas"`
	CreadoEn           time.Time `json:"creado_en" db:"creado_en"`
}

type RutinaEjercicio struct {
	ID        int       `json:"id" db:"id"`
	RutinaID  int       `json:"rutina_id" db:"rutina_id"`
	Ejercicio string    `json:"ejercicio" db:"ejercicio"`
	Orden     int       `json:"orden" db:"orden"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
}

type RutinaRequest struct {
	TipoEntrenamiento string   `json:"tipo_entrenamiento" binding:"required,max=30"`
	Nombre            string   `json:"nombre" binding:"required,min=2,max=150"`
	Nivel             string   `json:"nivel" binding:"required,oneof=Principiante Intermedio Avanzado"`
	Descripcion       string   `json:"descripcion"`
	DuracionSemanas   int      `json:"duracion_semanas" binding:"required,min=1"`
	Ejercicios        []string `json:"ejercicios"`
}

type RutinaUpdateRequest struct {
	TipoEntrenamiento *string `json:"tipo_entrenamiento" binding:"omitempty,max=30"`
	Nombre            *string `json:"nombre" binding:"omitempty,min=2,max=150"`
	Nivel             *string `json:"nivel" binding:"omitempty,oneof=Principiante Intermedio Avanzado"`
	Descripcion       *string `json:"descripcion" binding:"omitempty"`
	DuracionSemanas   *int    `json:"duracion_semanas" binding:"omitempty,min=1"`
	Ejercicios        []string `json:"ejercicios"`
}





type RutinaCompleta struct {
	Rutina
	Ejercicios []RutinaEjercicio `json:"ejercicios"`
}

const (
	NivelPrincipiante = "Principiante"
	NivelIntermedio   = "Intermedio"
	NivelAvanzado     = "Avanzado"
)



// Funciones constructoras de respuesta
func RespuestaExitosa(message string, data interface{}) RespuestaAPI {
	return RespuestaAPI{Success: true, Message: message, Data: data}
}

func RespuestaPaginadaExitosa(message string, data interface{}, pagina, porPagina int, total int64) RespuestaPaginada {
	return RespuestaPaginada{Success: true, Message: message, Data: data, Total: total, Pagina: pagina}
}

func RespuestaErrorGeneral(message string, err error) RespuestaError {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return RespuestaError{Success: false, Message: message, Error: errStr}
}