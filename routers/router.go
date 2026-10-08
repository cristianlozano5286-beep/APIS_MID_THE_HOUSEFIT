package routers

import (
	"api_mid_the_housefit/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/nutricion/calcular", &controllers.NutricionController{}, "post:CalcularIMC")
	beego.Router("/nutricion/historial", &controllers.NutricionController{}, "get:HistorialIMC")

	beego.Router("/nutricion/guias", &controllers.GuiaNutricionalController{}, "get:ListarGuias;post:CrearGuia")
	beego.Router("/nutricion/guias/:id", &controllers.GuiaNutricionalController{}, "delete:EliminarGuia")
}