package models

import "time"

type SolicitudLogin struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type SolicitudRegistro struct {
	Nombre   string `json:"nombre"`
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type SolicitudCorreo struct {
	Correo string `json:"correo"`
}

type SolicitudVerificar struct {
	Correo string `json:"correo"`
	Codigo string `json:"codigo"`
}

type SolicitudRestablecer struct {
	Correo        string `json:"correo"`
	Codigo        string `json:"codigo"`
	NuevaPassword string `json:"nuevaPassword"`
}

// el codigo se guarda en memoria mientras el MID este prendido
type CodigoVerificacion struct {
	Correo   string
	Codigo   string
	Creado   time.Time
	Vence    time.Time
	Intentos int
	Activo   bool
}
