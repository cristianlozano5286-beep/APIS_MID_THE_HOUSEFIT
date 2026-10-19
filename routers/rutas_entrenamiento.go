package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/controllers"
)

func init() {
	// ---------- clases y reservas ----------
	beego.Router("/api/v1/clases", &controllers.ClasesController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/clases/:id", &controllers.ClasesController{}, "delete:Borrar")
}
