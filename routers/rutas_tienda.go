package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/controllers"
)

func init() {
	// ---------- tienda ----------
	beego.Router("/api/v1/productos", &controllers.ProductosController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/productos/:id", &controllers.ProductosController{}, "put:Actualizar;delete:Borrar")
}
