package models

// GuiaNutricional, una por categoria de IMC.
type GuiaNutricional struct {
	ID              int    `json:"id"`
	CategoriaIMC    string `json:"categoria_imc"`
	Titulo          string `json:"titulo"`
	Recomendaciones string `json:"recomendaciones,omitempty"`
	EjemploComidas  string `json:"ejemplo_comidas,omitempty"`
}

// ResultadoIMC guarda cada calculo que hace un usuario (diagrama calcular_imc).
type ResultadoIMC struct {
	ID         int       `json:"id"`
	UsuarioCorreo string `json:"usuario_correo,omitempty"`
	PesoKg     float64   `json:"peso_kg"`
	EstaturaM  float64   `json:"estatura_m"`
	IMC        float64   `json:"imc"`
	Categoria  string    `json:"categoria"`
	Fecha      string    `json:"fecha"`
}
