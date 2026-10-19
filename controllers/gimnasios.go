package controllers

import (
	"sort"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"
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

// GET /api/v1/favoritos
// los gimnasios que marco la persona que inicio sesion
func (c *GimnasiosController) ListarFavoritos() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	var favoritos []models.Favorito
	if !listarColeccion(&c.Controller, "favoritos", &favoritos) {
		return
	}

	var gimnasios []models.Gimnasio
	if !listarColeccion(&c.Controller, "gimnasios", &gimnasios) {
		return
	}

	resultado := []models.Gimnasio{}
	for _, f := range favoritos {
		if strconv.Itoa(f.UsuarioId) != sesion.ID {
			continue
		}
		for _, g := range gimnasios {
			if g.Id == f.GimnasioId {
				resultado = append(resultado, g)
			}
		}
	}

	responder(&c.Controller, 200, resultado)
}

// POST /api/v1/gimnasios/:id/favorito
// si ya era favorito lo quita, si no lo agrega
func (c *GimnasiosController) AlternarFavorito() {
	sesion, valido := exigirSesion(&c.Controller)
	if !valido {
		return
	}

	id, valido := leerId(&c.Controller)
	if !valido {
		return
	}
	gimnasioId, _ := strconv.Atoi(id)
	usuarioId, _ := strconv.Atoi(sesion.ID)

	crud := services.NuevoServicioCrud()

	// el gimnasio tiene que existir
	var gimnasio models.Gimnasio
	err := crud.Obtener("gimnasios", id, &gimnasio)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	var favoritos []models.Favorito
	if !listarColeccion(&c.Controller, "favoritos", &favoritos) {
		return
	}

	for _, f := range favoritos {
		if f.UsuarioId == usuarioId && f.GimnasioId == gimnasioId {
			err = crud.Borrar("favoritos", strconv.Itoa(f.Id))
			if err != nil {
				errorCrud(&c.Controller, err)
				return
			}
			responder(&c.Controller, 200, map[string]interface{}{"esFavorito": false})
			return
		}
	}

	nuevo := models.Favorito{UsuarioId: usuarioId, GimnasioId: gimnasioId}
	err = crud.Crear("favoritos", nuevo, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, map[string]interface{}{"esFavorito": true})
}
