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

// POST /api/v1/gimnasios
func (c *GimnasiosController) Crear() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	var gimnasio models.Gimnasio

	if !leerCuerpo(&c.Controller, &gimnasio) {
		return
	}

	if strings.TrimSpace(gimnasio.Nombre) == "" || strings.TrimSpace(gimnasio.Direccion) == "" || strings.TrimSpace(gimnasio.Ciudad) == "" {
		responderError(&c.Controller, 400, "El nombre, la dirección y la ciudad son obligatorios.")
		return
	}

	crearEnColeccion(&c.Controller, "gimnasios", gimnasio)
}

// PUT /api/v1/gimnasios/:id
func (c *GimnasiosController) Actualizar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	actualizarEnColeccion(&c.Controller, "gimnasios")
}

// DELETE /api/v1/gimnasios/:id
func (c *GimnasiosController) Borrar() {
	_, valido := exigirRol(&c.Controller, "Administrador")
	if !valido {
		return
	}

	borrarDeColeccion(&c.Controller, "gimnasios")
}
