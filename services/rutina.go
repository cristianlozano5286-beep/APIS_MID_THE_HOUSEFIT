package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_housefit/models"
)

type RutinaService struct {
	db *sql.DB
}

func NewRutinaService() *RutinaService {
	return &RutinaService{db: GetDB()}
}

func (s *RutinaService) Create(req models.RutinaRequest) (*models.RutinaCompleta, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var rutina models.Rutina
	err = tx.QueryRow(`
		INSERT INTO rutinas (tipo_entrenamiento, nombre, nivel, descripcion, duracion_semanas)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, tipo_entrenamiento, nombre, nivel, descripcion, duracion_semanas, creado_en
	`, req.TipoEntrenamiento, req.Nombre, req.Nivel, req.Descripcion, req.DuracionSemanas).Scan(
		&rutina.ID, &rutina.TipoEntrenamiento, &rutina.Nombre, &rutina.Nivel, &rutina.Descripcion, &rutina.DuracionSemanas, &rutina.CreadoEn,
	)
	if err != nil {
		return nil, err
	}

	ejercicios := []models.RutinaEjercicio{}
	for i, ejercicio := range req.Ejercicios {
		var re models.RutinaEjercicio
		err = tx.QueryRow(`
			INSERT INTO rutina_ejercicios (rutina_id, ejercicio, orden)
			VALUES ($1, $2, $3)
			RETURNING id, rutina_id, ejercicio, orden, creado_en
		`, rutina.ID, ejercicio, i+1).Scan(&re.ID, &re.RutinaID, &re.Ejercicio, &re.Orden, &re.CreadoEn)
		if err != nil {
			return nil, err
		}
		ejercicios = append(ejercicios, re)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.RutinaCompleta{
		Rutina:     rutina,
		Ejercicios: ejercicios,
	}, nil
}

func (s *RutinaService) GetByID(id int) (*models.RutinaCompleta, error) {
	var rutina models.Rutina
	err := s.db.QueryRow(`
		SELECT id, tipo_entrenamiento, nombre, nivel, descripcion, duracion_semanas, creado_en
		FROM rutinas WHERE id = $1
	`, id).Scan(&rutina.ID, &rutina.TipoEntrenamiento, &rutina.Nombre, &rutina.Nivel, &rutina.Descripcion, &rutina.DuracionSemanas, &rutina.CreadoEn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`
		SELECT id, rutina_id, ejercicio, orden, creado_en
		FROM rutina_ejercicios WHERE rutina_id = $1 ORDER BY orden
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ejercicios := []models.RutinaEjercicio{}
	for rows.Next() {
		var re models.RutinaEjercicio
		err := rows.Scan(&re.ID, &re.RutinaID, &re.Ejercicio, &re.Orden, &re.CreadoEn)
		if err != nil {
			return nil, err
		}
		ejercicios = append(ejercicios, re)
	}

	return &models.RutinaCompleta{
		Rutina:     rutina,
		Ejercicios: ejercicios,
	}, nil
}

func (s *RutinaService) Update(id int, req models.RutinaUpdateRequest) (*models.RutinaCompleta, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	setClauses := []string{}
	values := []interface{}{}
	i := 1

	if req.TipoEntrenamiento != nil {
		setClauses = append(setClauses, fmt.Sprintf("tipo_entrenamiento = $%d", i))
		values = append(values, *req.TipoEntrenamiento)
		i++
	}
	if req.Nombre != nil {
		setClauses = append(setClauses, fmt.Sprintf("nombre = $%d", i))
		values = append(values, *req.Nombre)
		i++
	}
	if req.Nivel != nil {
		setClauses = append(setClauses, fmt.Sprintf("nivel = $%d", i))
		values = append(values, *req.Nivel)
		i++
	}
	if req.Descripcion != nil {
		setClauses = append(setClauses, fmt.Sprintf("descripcion = $%d", i))
		values = append(values, *req.Descripcion)
		i++
	}
	if req.DuracionSemanas != nil {
		setClauses = append(setClauses, fmt.Sprintf("duracion_semanas = $%d", i))
		values = append(values, *req.DuracionSemanas)
		i++
	}

	if len(setClauses) > 0 {
		values = append(values, id)
		query := fmt.Sprintf("UPDATE rutinas SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)
		_, err = tx.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	}

	if req.Ejercicios != nil {
		_, err = tx.Exec(`DELETE FROM rutina_ejercicios WHERE rutina_id = $1`, id)
		if err != nil {
			return nil, err
		}
		for i, ejercicio := range req.Ejercicios {
			_, err = tx.Exec(`
				INSERT INTO rutina_ejercicios (rutina_id, ejercicio, orden)
				VALUES ($1, $2, $3)
			`, id, ejercicio, i+1)
			if err != nil {
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *RutinaService) Delete(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM rutina_ejercicios WHERE rutina_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM rutinas WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *RutinaService) List(filters map[string]interface{}, pagina, porPagina int) ([]models.RutinaCompleta, int64, error) {
	whereClauses := []string{}
	values := []interface{}{}
	i := 1

	for col, val := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM rutinas %s", whereClause)
	err := s.db.QueryRow(countQuery, values...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	values = append(values, porPagina, offset)

	query := fmt.Sprintf(`
		SELECT id, tipo_entrenamiento, nombre, nivel, descripcion, duracion_semanas, creado_en
		FROM rutinas %s ORDER BY creado_en DESC LIMIT $%d OFFSET $%d
	`, whereClause, i, i+1)

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rutinas []models.RutinaCompleta
	for rows.Next() {
		var r models.Rutina
		err := rows.Scan(&r.ID, &r.TipoEntrenamiento, &r.Nombre, &r.Nivel, &r.Descripcion, &r.DuracionSemanas, &r.CreadoEn)
		if err != nil {
			return nil, 0, err
		}

		ejRows, err := s.db.Query(`
			SELECT id, rutina_id, ejercicio, orden, creado_en
			FROM rutina_ejercicios WHERE rutina_id = $1 ORDER BY orden
		`, r.ID)
		if err != nil {
			return nil, 0, err
		}

		ejercicios := []models.RutinaEjercicio{}
		for ejRows.Next() {
			var re models.RutinaEjercicio
			err := ejRows.Scan(&re.ID, &re.RutinaID, &re.Ejercicio, &re.Orden, &re.CreadoEn)
			if err != nil {
				ejRows.Close()
				return nil, 0, err
			}
			ejercicios = append(ejercicios, re)
		}
		ejRows.Close()

		rutinas = append(rutinas, models.RutinaCompleta{
			Rutina:     r,
			Ejercicios: ejercicios,
		})
	}

	return rutinas, total, nil
}