package controllers

import (
	"net/http"
	"sort"

	"api_mid_the_housefit/models"
)
func Metricas(w http.ResponseWriter, r *http.Request) {
	if _, ok := requerirRol(w, r, "Administrador"); !ok {
		return
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	porMes := map[string]float64{}
	for _, t := range db.transacciones {
		mes := t.Fecha
		if len(mes) >= 7 {
			mes = mes[:7]
		}
		porMes[mes] += t.Monto
	}
	ultimas := append([]models.Resena{}, db.resenas...)
	sort.Slice(ultimas, func(i, j int) bool { return ultimas[i].Fecha > ultimas[j].Fecha })
	if len(ultimas) > 3 {
		ultimas = ultimas[:3]
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"totales": map[string]int{
			"gimnasios": len(db.gimnasios), "instructores": len(db.instructores),
			"productos": len(db.productos), "usuarios": len(db.usuarios),
			"reservas": len(db.reservas), "transacciones": len(db.transacciones),
		},
		"pagos_por_mes":   porMes,
		"ultimas_resenas": ultimas,
	})
}










