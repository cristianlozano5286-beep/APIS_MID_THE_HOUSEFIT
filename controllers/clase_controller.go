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

// CrearClase responde a POST /clases 
func (c *ClaseController) CrearClase(w http.ResponseWriter, r *http.Request) {
	var in models.ClaseDisponible
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "payload JSON inválido", http.StatusBadRequest)
		return
	}

	if in.TipoClase == "" || in.GimnasioID == 0 {
		http.Error(w, "tipo_clase y gimnasio_id son obligatorios", http.StatusBadRequest)
		return
	}

	creada, err := c.service.Crear(r.Context(), in)
	if err != nil {
		http.Error(w, "error interno al crear la clase", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(creada)
}

// EliminarClase responde a DELETE /clases/{id} 
func (c *ClaseController) EliminarClase(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := c.service.Eliminar(r.Context(), id); err != nil {
		http.Error(w, "clase no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}