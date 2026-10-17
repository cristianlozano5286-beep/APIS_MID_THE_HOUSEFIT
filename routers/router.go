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
	beego.Router("/api/v1/resumen", &controllers.UsuariosController{}, "get:Resumen")

	// ---------- gimnasios y favoritos ----------
	beego.Router("/api/v1/gimnasios", &controllers.GimnasiosController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/gimnasios/ciudades", &controllers.GimnasiosController{}, "get:Ciudades")
	beego.Router("/api/v1/gimnasios/:id", &controllers.GimnasiosController{}, "get:Obtener;put:Actualizar;delete:Borrar")
	beego.Router("/api/v1/gimnasios/:id/favorito", &controllers.GimnasiosController{}, "post:AlternarFavorito")
	beego.Router("/api/v1/favoritos", &controllers.GimnasiosController{}, "get:ListarFavoritos")

	// ---------- instructores y tienda ----------
	beego.Router("/api/v1/instructores", &controllers.InstructoresController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/instructores/:id", &controllers.InstructoresController{}, "put:Actualizar;delete:Borrar")
	beego.Router("/api/v1/productos", &controllers.ProductosController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/productos/:id", &controllers.ProductosController{}, "put:Actualizar;delete:Borrar")

	// ---------- clases y reservas ----------
	beego.Router("/api/v1/clases", &controllers.ClasesController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/clases/:id", &controllers.ClasesController{}, "delete:Borrar")
	beego.Router("/api/v1/clases/:id/reservar", &controllers.ClasesController{}, "post:Reservar")
	beego.Router("/api/v1/reservas", &controllers.ClasesController{}, "get:MisReservas")
	beego.Router("/api/v1/reservas/:id/cancelar", &controllers.ClasesController{}, "put:CancelarReserva")

	// ---------- blog ----------
	beego.Router("/api/v1/resenas", &controllers.BlogController{}, "get:ListarResenas;post:CrearResena")
	beego.Router("/api/v1/resenas/promedio/:gimnasioId", &controllers.BlogController{}, "get:Promedio")
	beego.Router("/api/v1/resenas/:id", &controllers.BlogController{}, "delete:BorrarResena")
	beego.Router("/api/v1/noticias", &controllers.BlogController{}, "get:ListarNoticias;post:CrearNoticia")
	beego.Router("/api/v1/noticias/:id", &controllers.BlogController{}, "delete:BorrarNoticia")

	// ---------- rutinas, guias e IMC ----------
	beego.Router("/api/v1/rutinas", &controllers.ContenidoController{}, "get:ListarRutinas;post:CrearRutina")
	beego.Router("/api/v1/rutinas/:id", &controllers.ContenidoController{}, "delete:BorrarRutina")
	beego.Router("/api/v1/guias-fisicas", &controllers.ContenidoController{}, "get:ListarGuiasFisicas;post:CrearGuiaFisica")
	beego.Router("/api/v1/guias-fisicas/:id", &controllers.ContenidoController{}, "delete:BorrarGuiaFisica")
	beego.Router("/api/v1/guias-nutricionales", &controllers.ContenidoController{}, "get:ListarGuiasNutricionales;post:CrearGuiaNutricional")
	beego.Router("/api/v1/guias-nutricionales/:id", &controllers.ContenidoController{}, "delete:BorrarGuiaNutricional")
	beego.Router("/api/v1/imc", &controllers.ContenidoController{}, "get:CalcularIMC")

	// ---------- pagos ----------
	beego.Router("/api/v1/pagos/metodos", &controllers.PagosController{}, "get:ListarMetodos")
	beego.Router("/api/v1/pagos/transacciones", &controllers.PagosController{}, "get:MisTransacciones")
	beego.Router("/api/v1/pagos", &controllers.PagosController{}, "post:Pagar")
}
