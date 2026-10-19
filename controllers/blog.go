package controllers

import (
	"math"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
)

type BlogController struct {
	beego.Controller
}

// GET /api/v1/resenas?gimnasioId=1
func (c *BlogController) ListarResenas() {
	var resenas []models.Resena
	if !listarColeccion(&c.Controller, "resenas", &resenas) {
		return
	}

	gimnasioId := c.GetString("gimnasioId")

	resultado := []models.Resena{}
	for _, r := range resenas {
		if gimnasioId != "" && strconv.Itoa(r.GimnasioId) != gimnasioId {
			continue
		}
		resultado = append(resultado, r)
	}

	responder(&c.Controller, 200, resultado)
}

// GET /api/v1/resenas/promedio/:gimnasioId
func (c *BlogController) Promedio() {
	gimnasioId := c.Ctx.Input.Param(":gimnasioId")

	var resenas []models.Resena
	if !listarColeccion(&c.Controller, "resenas", &resenas) {
		return
	}

	suma := 0
	total := 0
	for _, r := range resenas {
		if strconv.Itoa(r.GimnasioId) == gimnasioId {
			suma += r.Puntuacion
			total++
		}
	}

	promedio := 0.0
	if total > 0 {
		promedio = math.Round(float64(suma)/float64(total)*10) / 10
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"promedio": promedio,
		"total":    total,
	})
}

// POST /api/v1/resenas
func (c *BlogController) CrearResena() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	var resena models.Resena

	if !leerCuerpo(&c.Controller, &resena) {
		return
	}

	if resena.GimnasioId == 0 || strings.TrimSpace(resena.Titulo) == "" || strings.TrimSpace(resena.Cuerpo) == "" {
		responderError(&c.Controller, 400, "El gimnasio, el título y el comentario son obligatorios.")
		return
	}

	if resena.Puntuacion < 1 || resena.Puntuacion > 5 {
		responderError(&c.Controller, 400, "La puntuación debe ser de 1 a 5.")
		return
	}

	// el autor y la fecha los pone el MID, no el cliente
	resena.Usuario = sesion.Nombre
	resena.Fecha = hoy()
	if resena.Imagen == "" {
		resena.Imagen = "⭐"
	}

	crearEnColeccion(&c.Controller, "resenas", resena)
}

// DELETE /api/v1/resenas/:id
func (c *BlogController) BorrarResena() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "resenas")
}

// GET /api/v1/noticias
func (c *BlogController) ListarNoticias() {
	var noticias []models.Noticia
	if !listarColeccion(&c.Controller, "noticias", &noticias) {
		return
	}

	responder(&c.Controller, 200, noticias)
}

// POST /api/v1/noticias
func (c *BlogController) CrearNoticia() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	var noticia models.Noticia

	if !leerCuerpo(&c.Controller, &noticia) {
		return
	}

	if strings.TrimSpace(noticia.Titulo) == "" || strings.TrimSpace(noticia.Cuerpo) == "" {
		responderError(&c.Controller, 400, "El título y el contenido son obligatorios.")
		return
	}

	noticia.Fecha = hoy()

	crearEnColeccion(&c.Controller, "noticias", noticia)
}

// DELETE /api/v1/noticias/:id
func (c *BlogController) BorrarNoticia() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "noticias")
}
