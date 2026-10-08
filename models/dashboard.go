package models

type DashboardMetricas struct {
	TotalUsuarios       int     `json:"total_usuarios"`
	TotalGimnasios      int     `json:"total_gimnasios"`
	TotalInstructores   int     `json:"total_instructores"`
	TotalProductos      int     `json:"total_productos"`
	TotalClases         int     `json:"total_clases"`
	TotalReservas       int     `json:"total_reservas"`
	TotalResenas        int     `json:"total_resenas"`
	TotalNoticias       int     `json:"total_noticias"`
	TotalRutinas        int     `json:"total_rutinas"`
	TotalTransacciones  int     `json:"total_transacciones"`
	IngresosTotales     float64 `json:"ingresos_totales"`
	IngresosMesActual   float64 `json:"ingresos_mes_actual"`
	UsuariosNuevosMes   int     `json:"usuarios_nuevos_mes"`
	ReservasMesActual   int     `json:"reservas_mes_actual"`
	ClasesPopulares     []ClasePopular `json:"clases_populares"`
	GimnasiosTop        []GimnasioTop  `json:"gimnasios_top"`
}

type ClasePopular struct {
	TipoClase     string `json:"tipo_clase"`
	TotalReservas int    `json:"total_reservas"`
}

type GimnasioTop struct {
	ID              int     `json:"id"`
	Nombre          string  `json:"nombre"`
	TotalReservas   int     `json:"total_reservas"`
	RatingPromedio  float64 `json:"rating_promedio"`
}