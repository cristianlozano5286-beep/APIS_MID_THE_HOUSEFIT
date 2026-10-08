package controllers

import (
	"net/http"

	"api_mid_the_housefit/models"
)

// Subcontrolador de servicios: vive aqui para que gimnasio_controller no crezca.

// ListarServicios del gimnasio.
func ListarServicios(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	out := []models.GimnasioServicio{}
	for _, s := range db.servicios {
		if s.GimnasioID == id {
			out = append(out, s)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

// AgregarServicio a un gimnasio (admin).
func AgregarServicio(w http.ResponseWriter, r *http.Request) {
	if _, ok := requerirRol(w, r, "Administrador"); !ok {
		return
	}
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	var in models.GimnasioServicio
	if !decodeJSON(w, r, &in) {
		return
	}
	if in.Servicio == "" {
		writeError(w, http.StatusBadRequest, "servicio es obligatorio")
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	in.ID = db.next("servicios")
	in.GimnasioID = id
	db.servicios = append(db.servicios, in)
	writeJSON(w, http.StatusCreated, in)
}

// EliminarServicio por id (admin).
func EliminarServicio(w http.ResponseWriter, r *http.Request) {
	if _, ok := requerirRol(w, r, "Administrador"); !ok {
		return
	}
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, s := range db.servicios {
		if s.ID == id {
			db.servicios = append(db.servicios[:i], db.servicios[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writeError(w, http.StatusNotFound, "servicio no encontrado")
}
