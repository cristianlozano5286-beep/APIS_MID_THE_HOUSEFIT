package models

type Producto struct {
	Id              int     `json:"id,omitempty"`
	Nombre          string  `json:"nombre"`
	Categoria       string  `json:"categoria"`
	Precio          float64 `json:"precio"`
	Stock           int     `json:"stock"`
	Descripcion     string  `json:"descripcion"`
	GuiaNutricional string  `json:"guiaNutricional"`
	Imagen          string  `json:"imagen"`
}
