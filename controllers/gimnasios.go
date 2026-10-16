package controllers

import (
	"sort"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"thehousefit_mid/models"
)

type GimnasiosController struct {
	beego.Controller
}

// GET /api/v1/gimnasios?ciudad=Yopal&activo=true
func (c *GimnasiosController) Listar() {
	var gimnasios []models.Gimnasio
	if !listarColeccion(&c.Controller, "gimnasios", &gimnasios) {
		return
	}

	ciudad := strings.TrimSpace(c.GetString("ciudad"))
	soloActivos := c.GetString("activo") == "true"

	resultado := []models.Gimnasio{}
	for _, g := range gimnasios {
		if ciudad != "" && !strings.EqualFold(g.Ciudad, ciudad) {
			continue
		}
		if soloActivos && !g.Activo {
			continue
		}
		resultado = append(resultado, g)
	}

	responder(&c.Controller, 200, resultado)
}

// GET /api/v1/gimnasios/ciudades
func (c *GimnasiosController) Ciudades() {
	var gimnasios []models.Gimnasio
	if !listarColeccion(&c.Controller, "gimnasios", &gimnasios) {
		return
	}

	vistas := map[string]bool{}
	ciudades := []string{}
	for _, g := range gimnasios {
		if !vistas[g.Ciudad] {
			vistas[g.Ciudad] = true
			ciudades = append(ciudades, g.Ciudad)
		}
	}
	sort.Strings(ciudades)

	responder(&c.Controller, 200, ciudades)
}

// GET /api/v1/gimnasios/:id
func (c *GimnasiosController) Obtener() {
	obtenerDeColeccion(&c.Controller, "gimnasios")
}
