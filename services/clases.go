package services
package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_house_fit/models"
)

type ClaseService struct {
	db *sql.DB
}

func NewClaseService() *ClaseService {
	return &ClaseService{db: GetDB()}
}

func (s *ClaseService) Create(req models.ClaseDisponibleRequest) (*models.ClaseDisponible, error) {
	var clase models.ClaseDisponible
	err := s.db.QueryRow(`
		INSERT INTO clases_disponibles (gimnasio_id, tipo_clase, instructor_nombre, instructor_id, hora, cupos, duracion, precio)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, gimnasio_id, tipo_clase, instructor_nombre, instructor_id, hora, cupos, duracion, precio, creado_en
	`, req.GimnasioID, req.TipoClase, req.InstructorNombre, req.InstructorID, req.Hora, req.Cupos, req.Duracion, req.Precio).Scan(
		&clase.ID, &clase.GimnasioID, &clase.TipoClase, &clase.InstructorNombre, &clase.InstructorID,
		&clase.Hora, &clase.Cupos, &clase.Duracion, &clase.Precio, &clase.CreadoEn,
	)
	return &clase, err
}

func (s *ClaseService) GetByID(id int) (*models.ClaseDisponible, error) {
	var clase models.ClaseDisponible
	err := s.db.QueryRow(`
		SELECT id, gimnasio_id, tipo_clase, instructor_nombre, instructor_id, hora, cupos, duracion, precio, creado_en
		FROM clases_disponibles WHERE id = $1
	`, id).Scan(
		&clase.ID, &clase.GimnasioID, &clase.TipoClase, &clase.InstructorNombre, &clase.InstructorID,
		&clase.Hora, &clase.Cupos, &clase.Duracion, &clase.Precio, &clase.CreadoEn,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &clase, err
}

func (s *ClaseService) Update(id int, req models.ClaseDisponibleUpdateRequest) (*models.ClaseDisponible, error) {
	setClauses := []string{}
	values := []interface{}{}
	i := 1

	if req.GimnasioID != nil {
		setClauses = append(setClauses, fmt.Sprintf("gimnasio_id = $%d", i))
		values = append(values, *req.GimnasioID)
		i++
	}
	if req.TipoClase != nil {
		setClauses = append(setClauses, fmt.Sprintf("tipo_clase = $%d", i))
		values = append(values, *req.TipoClase)
		i++
	}
	if req.InstructorNombre != nil {
		setClauses = append(setClauses, fmt.Sprintf("instructor_nombre = $%d", i))
		values = append(values, *req.InstructorNombre)
		i++
	}
	if req.InstructorID != nil {
		setClauses = append(setClauses, fmt.Sprintf("instructor_id = $%d", i))
		values = append(values, *req.InstructorID)
		i++
	}
	if req.Hora != nil {
		setClauses = append(setClauses, fmt.Sprintf("hora = $%d", i))
		values = append(values, *req.Hora)
		i++
	}
	if req.Cupos != nil {
		setClauses = append(setClauses, fmt.Sprintf("cupos = $%d", i))
		values = append(values, *req.Cupos)
		i++
	}
	if req.Duracion != nil {
		setClauses = append(setClauses, fmt.Sprintf("duracion = $%d", i))
		values = append(values, *req.Duracion)
		i++
	}
	if req.Precio != nil {
		setClauses = append(setClauses, fmt.Sprintf("precio = $%d", i))
		values = append(values, *req.Precio)
		i++
	}

	if len(setClauses) == 0 {
		return s.GetByID(id)
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE clases_disponibles SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)
	_, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *ClaseService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM clases_disponibles WHERE id = $1`, id)
	return err
}

func (s *ClaseService) List(filters map[string]interface{}, pagina, porPagina int) ([]models.ClaseDisponible, int64, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM clases_disponibles %s", whereClause)
	err := s.db.QueryRow(countQuery, values...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	values = append(values, porPagina, offset)

	query := fmt.Sprintf(`
		SELECT id, gimnasio_id, tipo_clase, instructor_nombre, instructor_id, hora, cupos, duracion, precio, creado_en
		FROM clases_disponibles %s ORDER BY creado_en DESC LIMIT $%d OFFSET $%d
	`, whereClause, i, i+1)

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clases []models.ClaseDisponible
	for rows.Next() {
		var c models.ClaseDisponible
		err := rows.Scan(&c.ID, &c.GimnasioID, &c.TipoClase, &c.InstructorNombre, &c.InstructorID, &c.Hora, &c.Cupos, &c.Duracion, &c.Precio, &c.CreadoEn)
		if err != nil {
			return nil, 0, err
		}
		clases = append(clases, c)
	}

	return clases, total, nil
}

func (s *ClaseService) GetByGimnasio(gimnasioID int) ([]models.ClaseDisponible, error) {
	rows, err := s.db.Query(`
		SELECT id, gimnasio_id, tipo_clase, instructor_nombre, instructor_id, hora, cupos, duracion, precio, creado_en
		FROM clases_disponibles WHERE gimnasio_id = $1 ORDER BY hora
	`, gimnasioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clases []models.ClaseDisponible
	for rows.Next() {
		var c models.ClaseDisponible
		err := rows.Scan(&c.ID, &c.GimnasioID, &c.TipoClase, &c.InstructorNombre, &c.InstructorID, &c.Hora, &c.Cupos, &c.Duracion, &c.Precio, &c.CreadoEn)
		if err != nil {
			return nil, err
		}
		clases = append(clases, c)
	}
	return clases, nil
}