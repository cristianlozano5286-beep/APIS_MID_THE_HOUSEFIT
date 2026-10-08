package routers

import (
	"api_mid_the_housefit/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)

	beego.Router("/api/clases", &controllers.ClaseController{}, "get:Listar")
	beego.Router("/api/clases/:id", &controllers.ClaseController{}, "get:Ver")
}
