package models

type Rutina struct {
	Id                int      `json:"id,omitempty"`
	TipoEntrenamiento string   `json:"tipoEntrenamiento"`
	Nombre            string   `json:"nombre"`
	Nivel             string   `json:"nivel"`
	Descripcion       string   `json:"descripcion"`
	Ejercicios        []string `json:"ejercicios"`
	DuracionSemanas   int      `json:"duracionSemanas"`
}

type GuiaFisica struct {
	Id               int    `json:"id,omitempty"`
	Titulo           string `json:"titulo"`
	Metodo           string `json:"metodo"`
	Contenido        string `json:"contenido"`
	FechaPublicacion string `json:"fechaPublicacion,omitempty"`
}

type GuiaNutricional struct {
	Id              int    `json:"id,omitempty"`
	CategoriaIMC    string `json:"categoriaIMC"`
	Titulo          string `json:"titulo"`
	Recomendaciones string `json:"recomendaciones"`
	EjemploComidas  string `json:"ejemploComidas"`
}
