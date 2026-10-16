package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

var ErrNoEncontrado = errors.New("no encontrado")
var ErrDatoInvalido = errors.New("dato invalido")
var ErrYaExiste = errors.New("ya existe")

// habla con el CRUD, que es el que guarda en la base de datos
type ServicioCrud struct {
	URL     string
	Cliente *http.Client
}

func NuevoServicioCrud() *ServicioCrud {
	return &ServicioCrud{
		URL:     beego.AppConfig.DefaultString("url_crud", "http://localhost:3000"),
		Cliente: &http.Client{Timeout: 10 * time.Second},
	}
}

// trae toda una coleccion
func (s *ServicioCrud) Listar(coleccion string, destino interface{}) error {
	return s.pedir("GET", coleccion, nil, destino)
}

func (s *ServicioCrud) Obtener(coleccion string, id string, destino interface{}) error {
	return s.pedir("GET", coleccion+"/"+id, nil, destino)
}

func (s *ServicioCrud) Crear(coleccion string, dato interface{}, destino interface{}) error {
	return s.pedir("POST", coleccion, dato, destino)
}

func (s *ServicioCrud) Actualizar(coleccion string, id string, dato interface{}, destino interface{}) error {
	return s.pedir("PUT", coleccion+"/"+id, dato, destino)
}

func (s *ServicioCrud) Borrar(coleccion string, id string) error {
	return s.pedir("DELETE", coleccion+"/"+id, nil, nil)
}

func (s *ServicioCrud) pedir(metodo string, ruta string, dato interface{}, destino interface{}) error {
	var cuerpo bytes.Buffer

	if dato != nil {
		err := json.NewEncoder(&cuerpo).Encode(dato)
		if err != nil {
			return fmt.Errorf("error armando el JSON: %w", err)
		}
	}

	peticion, err := http.NewRequest(metodo, s.URL+"/"+ruta, &cuerpo)
	if err != nil {
		return fmt.Errorf("error creando la peticion: %w", err)
	}
	peticion.Header.Set("Content-Type", "application/json")

	respuesta, err := s.Cliente.Do(peticion)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con el CRUD: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode == http.StatusNotFound {
		return ErrNoEncontrado
	}

	// la base rechazo el dato (por ejemplo una puntuacion de 9)
	if respuesta.StatusCode == http.StatusBadRequest {
		return ErrDatoInvalido
	}

	// ya existe (por ejemplo un correo repetido)
	if respuesta.StatusCode == http.StatusConflict {
		return ErrYaExiste
	}

	if respuesta.StatusCode < 200 || respuesta.StatusCode > 299 {
		return fmt.Errorf("el CRUD respondio con codigo %d", respuesta.StatusCode)
	}

	if destino == nil || respuesta.StatusCode == http.StatusNoContent {
		return nil
	}

	err = json.NewDecoder(respuesta.Body).Decode(destino)
	if err != nil {
		return fmt.Errorf("error leyendo la respuesta del CRUD: %w", err)
	}

	return nil
}
