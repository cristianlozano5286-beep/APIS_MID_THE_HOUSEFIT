package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/controllers"
)

func init() {
	// ---------- tienda ----------
	beego.Router("/api/v1/productos", &controllers.ProductosController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/productos/:id", &controllers.ProductosController{}, "put:Actualizar;delete:Borrar")

	// ---------- blog ----------
	beego.Router("/api/v1/resenas", &controllers.BlogController{}, "get:ListarResenas;post:CrearResena")
	beego.Router("/api/v1/resenas/promedio/:gimnasioId", &controllers.BlogController{}, "get:Promedio")
	beego.Router("/api/v1/resenas/:id", &controllers.BlogController{}, "delete:BorrarResena")
	beego.Router("/api/v1/noticias", &controllers.BlogController{}, "get:ListarNoticias;post:CrearNoticia")
	beego.Router("/api/v1/noticias/:id", &controllers.BlogController{}, "delete:BorrarNoticia")

	// ---------- pagos ----------
	beego.Router("/api/v1/pagos/metodos", &controllers.PagosController{}, "get:ListarMetodos")
	beego.Router("/api/v1/pagos/transacciones", &controllers.PagosController{}, "get:MisTransacciones")
	beego.Router("/api/v1/pagos", &controllers.PagosController{}, "post:Pagar")
}
