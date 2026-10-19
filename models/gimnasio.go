package models

type Gimnasio struct {
	Id          int      `json:"id,omitempty"`
	Nombre      string   `json:"nombre"`
	Direccion   string   `json:"direccion"`
	Barrio      string   `json:"barrio"`
	Ciudad      string   `json:"ciudad"`
	Pais        string   `json:"pais"`
	Telefono    string   `json:"telefono"`
	Horario     string   `json:"horario"`
	Descripcion string   `json:"descripcion"`
	Imagen      string   `json:"imagen"`
	Fotos       []string `json:"fotos"`
	Precio      float64  `json:"precio"`
	Moneda      string   `json:"moneda"`
	Servicios   []string `json:"servicios"`
	Activo      bool     `json:"activo"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
}

type Instructor struct {
	Id           int    `json:"id,omitempty"`
	Nombre       string `json:"nombre"`
	Correo       string `json:"correo"`
	Telefono     string `json:"telefono"`
	Especialidad string `json:"especialidad"`
	Experiencia  int    `json:"experiencia"`
	Biografia    string `json:"biografia"`
	Imagen       string `json:"imagen"`
	GimnasioId   *int   `json:"gimnasioId"`
	Activo       bool   `json:"activo"`
}

// gimnasio favorito de un usuario
type Favorito struct {
	Id         int `json:"id,omitempty"`
	UsuarioId  int `json:"usuarioId"`
	GimnasioId int `json:"gimnasioId"`
}
