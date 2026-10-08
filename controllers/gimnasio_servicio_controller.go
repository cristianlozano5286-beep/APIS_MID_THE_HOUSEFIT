package controllers

import (
	"net/http"

	"api_mid_the_housefit/models"
)

// ListarServicios obtiene todos los servicios ofrecidos por un gimnasio específico.
func ListarServicios(w http.ResponseWriter, r *http.Request) {
	// 1. Obtenemos el ID del gimnasio desde la URL
	gimnasioID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	// 2. Bloqueamos la base de datos en modo lectura para evitar problemas si alguien más está escribiendo
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 3. Filtramos únicamente los servicios que pertenecen a este gimnasio
	var serviciosDelGimnasio []models.GimnasioServicio
	
	for _, servicio := range db.servicios {
		if servicio.GimnasioID == gimnasioID {
			serviciosDelGimnasio = append(serviciosDelGimnasio, servicio)
		}
	}

	// 4. Respondemos al cliente con la lista encontrada
	writeJSON(w, http.StatusOK, serviciosDelGimnasio)
}

// AgregarServicio permite a un administrador añadir un nuevo servicio al catálogo de un gimnasio.
func AgregarServicio(w http.ResponseWriter, r *http.Request) {
	// 1. Verificamos que el usuario tenga permisos de Administrador
	if _, autorizado := requerirRol(w, r, "Administrador"); !autorizado {
		return
	}

	// 2. Obtenemos el ID del gimnasio al que le vamos a agregar el servicio
	gimnasioID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	// 3. Leemos los datos del nuevo servicio que vienen en el cuerpo de la petición
	var nuevoServicio models.GimnasioServicio
	if exito := decodeJSON(w, r, &nuevoServicio); !exito {
		return
	}

	// 4. Validamos la información vital
	if nuevoServicio.Servicio == "" {
		writeError(w, http.StatusBadRequest, "El nombre del servicio es obligatorio")
		return
	}

	// 5. Guardamos en la base de datos (bloqueando para escritura)
	db.mu.Lock()
	defer db.mu.Unlock()

	nuevoServicio.ID = db.next("servicios")
	nuevoServicio.GimnasioID = gimnasioID

	db.servicios = append(db.servicios, nuevoServicio)

	// 6. Confirmamos que se creó correctamente devolviendo el servicio con su nuevo ID
	writeJSON(w, http.StatusCreated, nuevoServicio)
}

// EliminarServicio quita un servicio del sistema usando su ID único.
func EliminarServicio(w http.ResponseWriter, r *http.Request) {
	// 1. Solo los administradores pueden borrar servicios
	if _, autorizado := requerirRol(w, r, "Administrador"); !autorizado {
		return
	}

	// 2. Obtenemos el ID del servicio que vamos a eliminar
	servicioID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	// 3. Bloqueamos la base de datos para modificarla de forma segura
	db.mu.Lock()
	defer db.mu.Unlock()

	// 4. Buscamos el servicio en nuestra lista
	for indice, servicio := range db.servicios {
		if servicio.ID == servicioID {
			// Lo encontramos: lo eliminamos de la lista uniendo la parte anterior y posterior al índice
			db.servicios = append(db.servicios[:indice], db.servicios[indice+1:]...)
			
			// Respondemos que la operación fue exitosa pero no hay contenido que devolver (204)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Si el ciclo termina y no retornó, el servicio nunca existió
	writeError(w, http.StatusNotFound, "El servicio que intentas eliminar no fue encontrado")
}