package models

type Clase struct {
	Id         int     `json:"id,omitempty"`
	GimnasioId int     `json:"gimnasioId"`
	TipoClase  string  `json:"tipoClase"`
	Instructor string  `json:"instructor"`
	Hora       string  `json:"hora"`
	Cupos      int     `json:"cupos"`
	Duracion   string  `json:"duracion"`
	Precio     float64 `json:"precio"`
}

type Reserva struct {
	Id       int    `json:"id,omitempty"`
	Usuario  string `json:"usuario"`
	ClaseId  int    `json:"claseId"`
	Fecha    string `json:"fecha"`
	Personas int    `json:"personas"`
	Estado   string `json:"estado"`
}

// lo que manda el cliente para reservar
type SolicitudReserva struct {
	Fecha    string `json:"fecha"`
	Personas int    `json:"personas"`
}
