package routers

import (
	"api_mid_the_housefit/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	beego.Router("/api/guias-nutricionales", &controllers.GuiaNutricionalController{}, "get:ListarGuias")
	beego.Router("/api/guias-nutricionales", &controllers.GuiaNutricionalController{}, "post:CrearGuia")
	beego.Router("/api/guias-nutricionales/:id", &controllers.GuiaNutricionalController{}, "put:ActualizarGuia")
	beego.Router("/api/guias-nutricionales/:id", &controllers.GuiaNutricionalController{}, "delete:EliminarGuia")
	beego.Router("/api/clases", &controllers.ClaseController{}, "get:Listar")
	beego.Router("/api/clases/:id", &controllers.ClaseController{}, "get:Ver")
	beego.Router("/api/guias-rutinas/:id", &controllers.RutinaController{}, "get:VerGuiaRutina")
	beego.Router("/api/guias-rutinas/:id", &controllers.RutinaController{}, "put:ActualizarGuiaRutina")
	beego.Router("/api/guias-rutinas/:id", &controllers.RutinaController{}, "delete:EliminarGuiaRutina")
	beego.Router("/api/guias-rutinas/:id", &controllers.RutinaController{}, "post:CrearGuiaRutina")
	

	
}
