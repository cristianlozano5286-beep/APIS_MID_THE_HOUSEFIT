package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_housefit/models"
	"api_mid_the_housefit/services"
)

// ClaseController agrupa los handlers e inyecta la capa de negocio
type ClaseController struct {
	service services.ClaseService
}

func NewClaseController(s services.ClaseService) *ClaseController {
	return &ClaseController{service: s}
}

// ListarClases responde a GET /clases?gimnasio_id=&tipo=
func (c *ClaseController) ListarClases(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	gimnasioID, _ := strconv.Atoi(q.Get("gimnasio_id"))
	tipo := q.Get("tipo")

	clases, err := c.service.ObtenerFiltradas(r.Context(), gimnasioID, tipo)
	if err != nil {
		http.Error(w, "error al consultar clases", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clases)
}

// ObtenerClase responde a GET /clases/{id}
func (c *ClaseController) ObtenerClase(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id")) // Go 1.22+ estándar o extractor de tu router
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	clase, err := c.service.ObtenerPorID(r.Context(), id)
	if err != nil {
		http.Error(w, "clase no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clase)
}

