package models

type Resena struct {
	Id         int    `json:"id,omitempty"`
	GimnasioId int    `json:"gimnasioId"`
	Usuario    string `json:"usuario"`
	Titulo     string `json:"titulo"`
	Cuerpo     string `json:"cuerpo"`
	Imagen     string `json:"imagen"`
	Puntuacion int    `json:"puntuacion"`
	Fecha      string `json:"fecha,omitempty"`
}

type Noticia struct {
	Id         int    `json:"id,omitempty"`
	Titulo     string `json:"titulo"`
	Cuerpo     string `json:"cuerpo"`
	Imagen     string `json:"imagen"`
	Tipo       string `json:"tipo"`
	GimnasioId *int   `json:"gimnasioId"`
	Fecha      string `json:"fecha,omitempty"`
}
