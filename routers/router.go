package routers

import (
	"api_mid_the_housefit/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	
	// ==================== GUÍAS NUTRICIONALES ROUTES ====================
	beego.Router("/api/guias-nutricionales", &controllers.GuiaNutricionalController{}, "get:ListarGuias")
	beego.Router("/api/guias-nutricionales", &controllers.GuiaNutricionalController{}, "post:CrearGuia")
	beego.Router("/api/guias-nutricionales/:id", &controllers.GuiaNutricionalController{}, "delete:EliminarGuia")
}