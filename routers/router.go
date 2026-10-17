package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"

	"thehousefit_mid/controllers"
)

func init() {
	// deja que el cliente Angular llame al MID desde otro puerto
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{beego.AppConfig.DefaultString("origen_cliente", "http://localhost:4200")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// ---------- acceso ----------
	beego.Router("/api/v1/auth/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/api/v1/auth/registro", &controllers.AuthController{}, "post:Registro")
	beego.Router("/api/v1/auth/codigo", &controllers.AuthController{}, "post:PedirCodigo")
	beego.Router("/api/v1/auth/verificar", &controllers.AuthController{}, "post:VerificarCodigo")
	beego.Router("/api/v1/auth/recuperar", &controllers.AuthController{}, "post:Recuperar")
	beego.Router("/api/v1/auth/restablecer", &controllers.AuthController{}, "post:Restablecer")
	beego.Router("/api/v1/auth/perfil", &controllers.AuthController{}, "get:Perfil")

	// ---------- usuarios y roles (solo Administrador) ----------
	beego.Router("/api/v1/usuarios", &controllers.UsuariosController{}, "get:Listar")
	beego.Router("/api/v1/usuarios/:id/rol", &controllers.UsuariosController{}, "put:CambiarRol")
}
