package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_mid_the_housefit/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// Declaración de la estructura del controlador de Guías Nutricionales
type GuiaNutricionalController struct {
	beego.Controller
}

// @Title ListarGuias
// @Description Lista las guías nutricionales, con filtro opcional por categoría
// @Param   categoria     query    string    false        "Bajo peso, Peso normal, Sobrepeso u Obesidad"
// @Success 200 {object} []models.GuiaNutricional
// @Failure 500 Error de base de datos
// @router /guias [get]
func (c *GuiaNutricionalController) ListarGuias() {

	// --------------------------------------------------------
	// 1. Obtener la categoría enviada en la URL (opcional)
	// --------------------------------------------------------
	categoria := c.GetString("categoria")

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
		`SELECT id, categoria_imc, titulo, recomendaciones, ejemplo_comidas
        FROM guias_nutricionales
        WHERE $1::text = '' OR categoria_imc = $1::text
        ORDER BY id`,
		categoria,
	)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"error":   "Error consultando las guías",
			"detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer filas.Close()

	// --------------------------------------------------------
	// 3. Leer las filas y armar la lista
	// --------------------------------------------------------
	guias := []models.GuiaNutricional{}
	for filas.Next() {
		var g models.GuiaNutricional
		err = filas.Scan(&g.ID, &g.CategoriaIMC, &g.Titulo, &g.Recomendaciones, &g.EjemploComidas)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			c.Data["json"] = map[string]interface{}{
				"error": "Error leyendo las guías",
			}
			c.ServeJSON()
			return
		}
		guias = append(guias, g)
	}

	// --------------------------------------------------------
	// 4. Devolver la lista
	// --------------------------------------------------------
	c.Data["json"] = guias
	c.ServeJSON()
}

// @Title CrearGuia
// @Description Crea una guía nutricional
// @Param   body     body    models.GuiaNutricional     true        "Datos de la guía"
// @Success 201 {object} models.GuiaNutricional
// @Failure 400 Datos inválidos
// @Failure 500 Error de base de datos
// @router /guias [post]
func (c *GuiaNutricionalController) CrearGuia() {

	// --------------------------------------------------------
	// 1. Leer el JSON que llega en el body
	// --------------------------------------------------------
	var in models.GuiaNutricional
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &in)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error": "JSON inválido",
		}
		c.ServeJSON()
		return
	}

	if in.CategoriaIMC == "" || in.Titulo == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error": "categoria_imc y titulo son obligatorios",
		}
		c.ServeJSON()
		return
	}

	// La tabla solo acepta estas 4 categorías
	if in.CategoriaIMC != "Bajo peso" && in.CategoriaIMC != "Peso normal" &&
		in.CategoriaIMC != "Sobrepeso" && in.CategoriaIMC != "Obesidad" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error": "categoria_imc debe ser: Bajo peso, Peso normal, Sobrepeso u Obesidad",
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 2. Guardar en guias_nutricionales
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

	err = conexion.QueryRow(
		`INSERT INTO guias_nutricionales (categoria_imc, titulo, recomendaciones, ejemplo_comidas)
        VALUES ($1, $2, $3, $4) RETURNING id`,
		in.CategoriaIMC, in.Titulo, in.Recomendaciones, in.EjemploComidas,
	).Scan(&in.ID)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"error":   "No se pudo guardar la guía",
			"detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 3. Devolver la guía creada
	// --------------------------------------------------------
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = in
	c.ServeJSON()
}

// @Title EliminarGuia
// @Description Elimina una guía nutricional por su ID
// @Param   id     path    int     true        "ID de la guía"
// @Success 204 Eliminada
// @Failure 400 ID inválido
// @Failure 404 Guía no encontrada
// @router /guias/:id [delete]
func (c *GuiaNutricionalController) EliminarGuia() {

	// --------------------------------------------------------
	// 1. Obtener el ID enviado en la URL
	// --------------------------------------------------------
	idString := c.Ctx.Input.Param(":id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error": "El ID debe ser numérico",
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 2. Eliminar de la base de datos
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

	resultado, err := conexion.Exec(`DELETE FROM guias_nutricionales WHERE id = $1`, id)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"error":   "No se pudo eliminar la guía",
			"detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// --------------------------------------------------------
	// 3. Revisar si existía la guía
	// --------------------------------------------------------
	filasBorradas, _ := resultado.RowsAffected()
	if filasBorradas == 0 {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = map[string]interface{}{
			"error": "Guía no encontrada",
		}
		c.ServeJSON()
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
}