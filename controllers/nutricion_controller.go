package controllers

import (
	"encoding/json"
	"math"
	"net/http"

	"api_mid_the_housefit/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// Declaración de la estructura del controlador de Nutrición
type NutricionController struct {
	beego.Controller
}

// Datos que llegan en el body de la petición
type IMCEntrada struct {
	PesoKg        float64 `json:"peso_kg"`
	EstaturaM     float64 `json:"estatura_m"`
	UsuarioCorreo string  `json:"usuario_correo"`
}

// @Title CalcularIMC
// @Description Calcula el IMC, lo clasifica, lo guarda en la base de datos y devuelve la guía nutricional
// @Param   body     body    controllers.IMCEntrada     true        "Peso, estatura y correo"
// @Success 200 {object} models.ResultadoIMC
// @Failure 400 Datos inválidos
// @Failure 500 Error de base de datos
// @router /calcular [post]
func (c *NutricionController) CalcularIMC() {

	// --------------------------------------------------------
	// 1. Leer el JSON que llega en el body
	// --------------------------------------------------------
	var in IMCEntrada
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &in)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error": "JSON inválido",
		}
		c.ServeJSON()
		return
	}

	if in.PesoKg <= 0 || in.EstaturaM <= 0 {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error": "peso_kg y estatura_m deben ser mayores a 0",
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 2. Calcular el IMC y clasificarlo
	// --------------------------------------------------------
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

	// --------------------------------------------------------
	// 3. Obtener la conexión a la base de datos
	// --------------------------------------------------------
	conexion, err := orm.GetDB()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"error": "No fue posible conectarse a la base de datos",
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 4. Guardar el resultado en resultados_imc
	// --------------------------------------------------------
	// Si no envían correo se guarda NULL
	var correo interface{}
	if in.UsuarioCorreo != "" {
		correo = in.UsuarioCorreo
	}

	var id int
	err = conexion.QueryRow(
		`INSERT INTO resultados_imc (usuario_correo, peso_kg, estatura_m, imc, categoria)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		correo, in.PesoKg, in.EstaturaM, imc, categoria,
	).Scan(&id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error":   "No se pudo guardar el resultado (revise que el correo exista en usuarios)",
			"detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 5. Buscar la guía nutricional de esa categoría
	// --------------------------------------------------------
	var ptrGuia map[string]interface{} = nil

	var guiaID int
	var titulo, recomendaciones, ejemploComidas string
	err = conexion.QueryRow(
		`SELECT id, titulo, recomendaciones, ejemplo_comidas
		FROM guias_nutricionales WHERE categoria_imc = $1 LIMIT 1`,
		categoria,
	).Scan(&guiaID, &titulo, &recomendaciones, &ejemploComidas)

	// Si hay guía para la categoría se agrega, si no queda null
	if err == nil {
		ptrGuia = map[string]interface{}{
			"id":              guiaID,
			"categoria_imc":   categoria,
			"titulo":          titulo,
			"recomendaciones": recomendaciones,
			"ejemplo_comidas": ejemploComidas,
		}
	}

	// --------------------------------------------------------
	// 6. Devolver el resultado
	// --------------------------------------------------------
	c.Data["json"] = map[string]interface{}{
		"id":        id,
		"imc":       imc,
		"categoria": categoria,
		"guia":      ptrGuia, // Contendrá la guía o null si no existe
	}
	c.ServeJSON()
}

// @Title HistorialIMC
// @Description Devuelve los cálculos guardados, con filtro opcional por correo
// @Param   usuario     query    string     false        "Correo del usuario"
// @Success 200 {object} []models.ResultadoIMC
// @Failure 500 Error de base de datos
// @router /historial [get]
func (c *NutricionController) HistorialIMC() {

	// --------------------------------------------------------
	// 1. Obtener el correo enviado en la URL (opcional)
	// --------------------------------------------------------
	correo := c.GetString("usuario")

	// --------------------------------------------------------
	// 2. Consultar la base de datos
	// --------------------------------------------------------
	conexion, err := orm.GetDB()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"error": "No fue posible conectarse a la base de datos",
		}
		c.ServeJSON()
		return
	}

	filas, err := conexion.Query(
		`SELECT id, COALESCE(usuario_correo, ''), peso_kg, estatura_m, imc, categoria, fecha::text
		 FROM resultados_imc
		 WHERE $1::text = '' OR usuario_correo = $1::text
		 ORDER BY id DESC`,
		correo,
	)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"error":   "Error consultando el historial",
			"detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer filas.Close()

	// --------------------------------------------------------
	// 3. Leer las filas y armar la lista
	// --------------------------------------------------------
	resultados := []models.ResultadoIMC{}
	for filas.Next() {
		var res models.ResultadoIMC
		err = filas.Scan(&res.ID, &res.UsuarioCorreo, &res.PesoKg,
			&res.EstaturaM, &res.IMC, &res.Categoria, &res.Fecha)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			c.Data["json"] = map[string]interface{}{
				"error": "Error leyendo los resultados",
			}
			c.ServeJSON()
			return
		}
		resultados = append(resultados, res)
	}

	// --------------------------------------------------------
	// 4. Devolver la lista
	// --------------------------------------------------------
	c.Data["json"] = resultados
	c.ServeJSON()
}