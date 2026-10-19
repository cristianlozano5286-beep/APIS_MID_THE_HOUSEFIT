package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/controllers"
)

func init() {
	// ---------- gimnasios, favoritos e instructores ----------
	beego.Router("/api/v1/gimnasios", &controllers.GimnasiosController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/gimnasios/ciudades", &controllers.GimnasiosController{}, "get:Ciudades")
	beego.Router("/api/v1/gimnasios/:id", &controllers.GimnasiosController{}, "get:Obtener;put:Actualizar;delete:Borrar")
}
