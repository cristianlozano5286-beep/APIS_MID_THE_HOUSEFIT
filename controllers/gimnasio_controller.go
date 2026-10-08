package controllers

import (
	"net/http"
	"sort"
	"time"

	"api_mid_the_housefit/models"
)

// CiudadesDisponibles obtiene una lista de todas las ciudades donde hay gimnasios, sin repetirlas, ideal para el buscador "Dónde".
func CiudadesDisponibles(w http.ResponseWriter, r *http.Request) {
	// 1. Bloqueamos la base de datos en modo lectura
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 2. Usamos un mapa para registrar qué ciudades ya hemos visto y evitar duplicados
	ciudadesUnicas := map[string]bool{}
	for _, gimnasio := range db.gimnasios {
		ciudadesUnicas[gimnasio.Ciudad] = true
	}

	// 3. Convertimos el mapa en una lista normal de textos (strings)
	var listaCiudades []string
	for ciudad := range ciudadesUnicas {
		listaCiudades = append(listaCiudades, ciudad)
	}

	// 4. Ordenamos las ciudades alfabéticamente para que se vean mejor en la interfaz
	sort.Strings(listaCiudades)

	// 5. Enviamos la lista al cliente
	writeJSON(w, http.StatusOK, listaCiudades)
}

// ListarGimnasios devuelve el catálogo completo de gimnasios y permite filtrar por ciudad o por un servicio específico.
func ListarGimnasios(w http.ResponseWriter, r *http.Request) {
	// 1. Obtenemos los parámetros de búsqueda de la URL (ej: ?ciudad=Bogotá&servicio=Pesas)
	parametrosUrl := r.URL.Query()
	ciudadBuscada := parametrosUrl.Get("ciudad")
	servicioBuscado := parametrosUrl.Get("servicio")

	// 2. Bloqueamos en modo lectura
	db.mu.RLock()
	defer db.mu.RUnlock()

	var gimnasiosFiltrados []models.Gimnasio

	// 3. Revisamos cada gimnasio para ver si cumple con los filtros
	for _, gimnasio := range db.gimnasios {
		// Ignoramos los gimnasios que están inactivos
		if !gimnasio.Activo {
			continue
		}
		
		// Si el usuario filtró por ciudad y este gimnasio no está en esa ciudad, lo saltamos
		if ciudadBuscada != "" && gimnasio.Ciudad != ciudadBuscada {
			continue
		}
		
		// Si el usuario busca un servicio y el gimnasio no lo tiene, lo saltamos
		if servicioBuscado != "" && !tieneServicio(gimnasio.ID, servicioBuscado) {
			continue
		}
		
		// Si pasó todos los filtros, lo agregamos a nuestra lista final
		gimnasiosFiltrados = append(gimnasiosFiltrados, gimnasio)
	}

	// 4. Respondemos con la lista de gimnasios que pasaron la prueba
	writeJSON(w, http.StatusOK, gimnasiosFiltrados)
}

// tieneServicio es una función auxiliar para saber si un gimnasio ofrece un servicio en particular.
func tieneServicio(gimnasioID int, servicioBuscado string) bool {
	for _, servicioActual := range db.servicios {
		if servicioActual.GimnasioID == gimnasioID && servicioActual.Servicio == servicioBuscado {
			return true
		}
	}
	return false
}

// GimnasioDetalle representa la ficha completa de un gimnasio, agrupando su información básica, sus fotos y sus servicios.
type GimnasioDetalle struct {
	models.Gimnasio
	Fotos     []models.GimnasioFoto     `json:"fotos"`
	Servicios []models.GimnasioServicio `json:"servicios"`
}

// VerGimnasio muestra la ficha detallada de un solo gimnasio, incluyendo sus fotos y los servicios que ofrece.
func VerGimnasio(w http.ResponseWriter, r *http.Request) {
	// 1. Obtenemos el ID del gimnasio a consultar
	gimnasioID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	// 2. Bloqueamos la base de datos en modo lectura
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 3. Buscamos el gimnasio por su ID
	for _, gimnasio := range db.gimnasios {
		if gimnasio.ID == gimnasioID {
			// Preparamos la estructura detallada
			detalle := GimnasioDetalle{
				Gimnasio:  gimnasio,
				Fotos:     []models.GimnasioFoto{},
				Servicios: []models.GimnasioServicio{},
			}

			// Buscamos las fotos que le pertenecen a este gimnasio
			for _, foto := range db.fotos {
				if foto.GimnasioID == gimnasioID {
					detalle.Fotos = append(detalle.Fotos, foto)
				}
			}

			// Buscamos los servicios que ofrece este gimnasio
			for _, servicio := range db.servicios {
				if servicio.GimnasioID == gimnasioID {
					detalle.Servicios = append(detalle.Servicios, servicio)
				}
			}

			// Enviamos la ficha completa al cliente
			writeJSON(w, http.StatusOK, detalle)
			return
		}
	}

	// Si terminamos de buscar y no lo encontramos:
	writeError(w, http.StatusNotFound, "No pudimos encontrar el gimnasio solicitado")
}

// CrearGimnasio permite a un administrador registrar un nuevo gimnasio en la plataforma. Por defecto nace "Activo".
func CrearGimnasio(w http.ResponseWriter, r *http.Request) {
	// 1. Validamos que solo los administradores puedan hacer esto
	if _, autorizado := requerirRol(w, r, "Administrador"); !autorizado {
		return
	}

	// 2. Extraemos los datos del nuevo gimnasio enviados en la petición
	var nuevoGimnasio models.Gimnasio
	if exito := decodeJSON(w, r, &nuevoGimnasio); !exito {
		return
	}

	// 3. Validamos que nos hayan enviado la información mínima requerida
	if nuevoGimnasio.Nombre == "" || nuevoGimnasio.Ciudad == "" {
		writeError(w, http.StatusBadRequest, "El nombre y la ciudad del gimnasio son datos obligatorios")
		return
	}

	// 4. Bloqueamos la base de datos para guardar de forma segura
	db.mu.Lock()
	defer db.mu.Unlock()

	// 5. Asignamos valores automáticos (ID, estado, fecha de creación y país por defecto)
	nuevoGimnasio.ID = db.next("gimnasios")
	nuevoGimnasio.Activo = true
	nuevoGimnasio.CreadoEn = time.Now()
	
	if nuevoGimnasio.Pais == "" {
		nuevoGimnasio.Pais = "Colombia"
	}

	// 6. Guardamos el gimnasio en nuestra lista
	db.gimnasios = append(db.gimnasios, nuevoGimnasio)

	// 7. Confirmamos la creación exitosa
	writeJSON(w, http.StatusCreated, nuevoGimnasio)
}

// ActualizarGimnasio permite a un administrador modificar los datos de un gimnasio existente.
func ActualizarGimnasio(w http.ResponseWriter, r *http.Request) {
	// 1. Verificamos permisos
	if _, autorizado := requerirRol(w, r, "Administrador"); !autorizado {
		return
	}

	// 2. Obtenemos el ID del gimnasio a editar
	gimnasioID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	// 3. Leemos los datos que el usuario quiere actualizar
	var datosActualizados models.Gimnasio
	if exito := decodeJSON(w, r, &datosActualizados); !exito {
		return
	}

	// 4. Bloqueamos la base de datos para modificarla
	db.mu.Lock()
	defer db.mu.Unlock()

	// 5. Buscamos el gimnasio y actualizamos solo los campos que fueron enviados (no están vacíos)
	for indice, gimnasioActual := range db.gimnasios {
		if gimnasioActual.ID == gimnasioID {
			
			if datosActualizados.Nombre != "" {
				db.gimnasios[indice].Nombre = datosActualizados.Nombre
			}
			if datosActualizados.Direccion != "" {
				db.gimnasios[indice].Direccion = datosActualizados.Direccion
			}
			if datosActualizados.Barrio != "" {
				db.gimnasios[indice].Barrio = datosActualizados.Barrio
			}
			if datosActualizados.Ciudad != "" {
				db.gimnasios[indice].Ciudad = datosActualizados.Ciudad
			}
			if datosActualizados.Telefono != "" {
				db.gimnasios[indice].Telefono = datosActualizados.Telefono
			}
			if datosActualizados.Horario != "" {
				db.gimnasios[indice].Horario = datosActualizados.Horario
			}
			if datosActualizados.Descripcion != "" {
				db.gimnasios[indice].Descripcion = datosActualizados.Descripcion
			}
			if datosActualizados.Precio != 0 {
				db.gimnasios[indice].Precio = datosActualizados.Precio
			}
			
			// Devolvemos el gimnasio con sus datos ya actualizados
			writeJSON(w, http.StatusOK, db.gimnasios[indice])
			return
		}
	}

	writeError(w, http.StatusNotFound, "No pudimos encontrar el gimnasio para actualizarlo")
}

// EliminarGimnasio permite a un administrador borrar un gimnasio del sistema por completo.
func EliminarGimnasio(w http.ResponseWriter, r *http.Request) {
	// 1. Verificamos permisos
	if _, autorizado := requerirRol(w, r, "Administrador"); !autorizado {
		return
	}

	// 2. Obtenemos el ID a eliminar
	gimnasioID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	// 3. Bloqueamos la base de datos
	db.mu.Lock()
	defer db.mu.Unlock()

	// 4. Buscamos y eliminamos el gimnasio
	for indice, gimnasio := range db.gimnasios {
		if gimnasio.ID == gimnasioID {
			// Cortamos la porción de la lista antes y después del gimnasio, excluyéndolo
			db.gimnasios = append(db.gimnasios[:indice], db.gimnasios[indice+1:]...)
			
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	writeError(w, http.StatusNotFound, "El gimnasio que intentas eliminar no existe")
}