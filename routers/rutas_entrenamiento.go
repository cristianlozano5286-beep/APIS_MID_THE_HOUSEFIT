package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/controllers"
)

func init() {
	// ---------- clases y reservas ----------
	beego.Router("/api/v1/clases", &controllers.ClasesController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/clases/:id", &controllers.ClasesController{}, "delete:Borrar")
	beego.Router("/api/v1/clases/:id/reservar", &controllers.ClasesController{}, "post:Reservar")
	beego.Router("/api/v1/reservas", &controllers.ClasesController{}, "get:MisReservas")
	beego.Router("/api/v1/reservas/:id/cancelar", &controllers.ClasesController{}, "put:CancelarReserva")

	// ---------- rutinas, guias e IMC ----------
	beego.Router("/api/v1/rutinas", &controllers.ContenidoController{}, "get:ListarRutinas;post:CrearRutina")
	beego.Router("/api/v1/rutinas/:id", &controllers.ContenidoController{}, "delete:BorrarRutina")
	beego.Router("/api/v1/guias-fisicas", &controllers.ContenidoController{}, "get:ListarGuiasFisicas;post:CrearGuiaFisica")
	beego.Router("/api/v1/guias-fisicas/:id", &controllers.ContenidoController{}, "delete:BorrarGuiaFisica")
	beego.Router("/api/v1/guias-nutricionales", &controllers.ContenidoController{}, "get:ListarGuiasNutricionales;post:CrearGuiaNutricional")
	beego.Router("/api/v1/guias-nutricionales/:id", &controllers.ContenidoController{}, "delete:BorrarGuiaNutricional")
}
