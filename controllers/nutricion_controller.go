package controllers

import (
	"math"
	"net/http"
	"time"

	"api_mid_the_housefit/models"
)

// CalcularIMC: calcula el IMC, lo clasifica, GUARDA el resultado
// y devuelve la guia nutricional de esa categoria.
func CalcularIMC(w http.ResponseWriter, r *http.Request) {
	var in struct {
		PesoKg     float64 `json:"peso_kg"`
		EstaturaM  float64 `json:"estatura_m"`
		UsuarioCorreo string `json:"usuario_correo"`
	}
	if !decodeJSON(w, r, &in) {
		return
	}
	if in.PesoKg <= 0 || in.EstaturaM <= 0 {
		writeError(w, http.StatusBadRequest, "peso_kg y estatura_m deben ser mayores a 0")
		return
	}
	imc := math.Round(in.PesoKg/(in.EstaturaM*in.EstaturaM)*10) / 10
	categoria := "Obesidad"
	switch {
	case imc < 18.5:
		categoria = "Bajo peso"
	case imc < 25:
		categoria = "Peso normal"
	case imc < 30:
		categoria = "Sobrepeso"
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	res := models.ResultadoIMC{
		ID: db.next("resultados_imc"), UsuarioCorreo: in.UsuarioCorreo,
		PesoKg: in.PesoKg, EstaturaM: in.EstaturaM,
		IMC: imc, Categoria: categoria, Fecha: time.Now().Format("2006-01-02"),
	}
	db.resultadosIMC = append(db.resultadosIMC, res)
	var guia *models.GuiaNutricional
	for i, g := range db.guiasNutricionales {
		if g.CategoriaIMC == categoria {
			guia = &db.guiasNutricionales[i]
			break
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"imc": imc, "categoria": categoria, "guia": guia,
	})
}

// HistorialIMC: calculos guardados, con filtro ?usuario= (correo).
func HistorialIMC(w http.ResponseWriter, r *http.Request) {
	correo := r.URL.Query().Get("usuario")
	db.mu.RLock()
	defer db.mu.RUnlock()
	out := []models.ResultadoIMC{}
	for _, res := range db.resultadosIMC {
		if correo == "" || res.UsuarioCorreo == correo {
			out = append(out, res)
		}
	}
	writeJSON(w, http.StatusOK, out)
}
