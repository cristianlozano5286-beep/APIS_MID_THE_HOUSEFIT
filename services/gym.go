package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_house_fit/models"
)

type GimnasioService struct {
	db *sql.DB
}

func NewGimnasioService() *GimnasioService {
	return &GimnasioService{db: GetDB()}
}

func (s *GimnasioService) Create(req models.GimnasioRequest) (*models.GimnasioCompleto, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var gimnasio models.Gimnasio
	err = tx.QueryRow(`
		INSERT INTO gimnasios (nombre, direccion, barrio, ciudad, pais, telefono, horario, descripcion, imagen, precio, moneda, activo, lat, lng)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, nombre, direccion, barrio, ciudad, pais, telefono, horario, descripcion, imagen, precio, moneda, activo, lat, lng, creado_en
	`, req.Nombre, req.Direccion, req.Barrio, req.Ciudad, req.Pais, req.Telefono, req.Horario,
		req.Descripcion, req.Imagen, req.Precio, req.Moneda, req.Activo, req.Lat, req.Lng).Scan(
		&gimnasio.ID, &gimnasio.Nombre, &gimnasio.Direccion, &gimnasio.Barrio, &gimnasio.Ciudad,
		&gimnasio.Pais, &gimnasio.Telefono, &gimnasio.Horario, &gimnasio.Descripcion, &gimnasio.Imagen,
		&gimnasio.Precio, &gimnasio.Moneda, &gimnasio.Activo, &gimnasio.Lat, &gimnasio.Lng, &gimnasio.CreadoEn,
	)
	if err != nil {
		return nil, err
	}

	for i, url := range req.Fotos {
		_, err = tx.Exec(`
			INSERT INTO gimnasio_fotos (gimnasio_id, url, orden)
			VALUES ($1, $2, $3)
		`, gimnasio.ID, url, i+1)
		if err != nil {
			return nil, err
		}
	}

	for _, servicio := range req.Servicios {
		_, err = tx.Exec(`
			INSERT INTO gimnasio_servicios (gimnasio_id, servicio)
			VALUES ($1, $2)
		`, gimnasio.ID, servicio)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return s.GetCompleto(gimnasio.ID)
}

func (s *GimnasioService) GetCompleto(id int) (*models.GimnasioCompleto, error) {
	var g models.Gimnasio
	err := s.db.QueryRow(`
		SELECT id, nombre, direccion, barrio, ciudad, pais, telefono, horario, descripcion, imagen, precio, moneda, activo, lat, lng, creado_en
		FROM gimnasios WHERE id = $1
	`, id).Scan(
		&g.ID, &g.Nombre, &g.Direccion, &g.Barrio, &g.Ciudad, &g.Pais, &g.Telefono, &g.Horario,
		&g.Descripcion, &g.Imagen, &g.Precio, &g.Moneda, &g.Activo, &g.Lat, &g.Lng, &g.CreadoEn,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	fotos, err := s.getFotos(id)
	if err != nil {
		return nil, err
	}

	servicios, err := s.getServicios(id)
	if err != nil {
		return nil, err
	}

	return &models.GimnasioCompleto{
		Gimnasio:  g,
		Fotos:     fotos,
		Servicios: servicios,
	}, nil
}

func (s *GimnasioService) getFotos(gimnasioID int) ([]models.GimnasioFoto, error) {
	rows, err := s.db.Query(`
		SELECT id, gimnasio_id, url, orden
		FROM gimnasio_fotos WHERE gimnasio_id = $1 ORDER BY orden
	`, gimnasioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fotos []models.GimnasioFoto
	for rows.Next() {
		var f models.GimnasioFoto
		err := rows.Scan(&f.ID, &f.GimnasioID, &f.URL, &f.Orden)
		if err != nil {
			return nil, err
		}
		fotos = append(fotos, f)
	}
	return fotos, nil
}

func (s *GimnasioService) getServicios(gimnasioID int) ([]models.GimnasioServicio, error) {
	rows, err := s.db.Query(`
		SELECT id, gimnasio_id, servicio
		FROM gimnasio_servicios WHERE gimnasio_id = $1
	`, gimnasioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servicios []models.GimnasioServicio
	for rows.Next() {
		var srv models.GimnasioServicio
		err := rows.Scan(&srv.ID, &srv.GimnasioID, &srv.Servicio)
		if err != nil {
			return nil, err
		}
		servicios = append(servicios, srv)
	}
	return servicios, nil
}

func (s *GimnasioService) Update(id int, req models.GimnasioUpdateRequest) (*models.GimnasioCompleto, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	setClauses := []string{}
	values := []interface{}{}
	i := 1

	if req.Nombre != nil {
		setClauses = append(setClauses, fmt.Sprintf("nombre = $%d", i))
		values = append(values, *req.Nombre)
		i++
	}
	if req.Direccion != nil {
		setClauses = append(setClauses, fmt.Sprintf("direccion = $%d", i))
		values = append(values, *req.Direccion)
		i++
	}
	if req.Barrio != nil {
		setClauses = append(setClauses, fmt.Sprintf("barrio = $%d", i))
		values = append(values, *req.Barrio)
		i++
	}
	if req.Ciudad != nil {
		setClauses = append(setClauses, fmt.Sprintf("ciudad = $%d", i))
		values = append(values, *req.Ciudad)
		i++
	}
	if req.Pais != nil {
		setClauses = append(setClauses, fmt.Sprintf("pais = $%d", i))
		values = append(values, *req.Pais)
		i++
	}
	if req.Telefono != nil {
		setClauses = append(setClauses, fmt.Sprintf("telefono = $%d", i))
		values = append(values, *req.Telefono)
		i++
	}
	if req.Horario != nil {
		setClauses = append(setClauses, fmt.Sprintf("horario = $%d", i))
		values = append(values, *req.Horario)
		i++
	}
	if req.Descripcion != nil {
		setClauses = append(setClauses, fmt.Sprintf("descripcion = $%d", i))
		values = append(values, *req.Descripcion)
		i++
	}
	if req.Imagen != nil {
		setClauses = append(setClauses, fmt.Sprintf("imagen = $%d", i))
		values = append(values, *req.Imagen)
		i++
	}
	if req.Precio != nil {
		setClauses = append(setClauses, fmt.Sprintf("precio = $%d", i))
		values = append(values, *req.Precio)
		i++
	}
	if req.Moneda != nil {
		setClauses = append(setClauses, fmt.Sprintf("moneda = $%d", i))
		values = append(values, *req.Moneda)
		i++
	}
	if req.Activo != nil {
		setClauses = append(setClauses, fmt.Sprintf("activo = $%d", i))
		values = append(values, *req.Activo)
		i++
	}
	if req.Lat != nil {
		setClauses = append(setClauses, fmt.Sprintf("lat = $%d", i))
		values = append(values, *req.Lat)
		i++
	}
	if req.Lng != nil {
		setClauses = append(setClauses, fmt.Sprintf("lng = $%d", i))
		values = append(values, *req.Lng)
		i++
	}

	if len(setClauses) > 0 {
		values = append(values, id)
		query := fmt.Sprintf("UPDATE gimnasios SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)
		_, err = tx.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	}

	if req.Fotos != nil {
		_, err = tx.Exec(`DELETE FROM gimnasio_fotos WHERE gimnasio_id = $1`, id)
		if err != nil {
			return nil, err
		}
		for i, url := range req.Fotos {
			_, err = tx.Exec(`
				INSERT INTO gimnasio_fotos (gimnasio_id, url, orden)
				VALUES ($1, $2, $3)
			`, id, url, i+1)
			if err != nil {
				return nil, err
			}
		}
	}

	if req.Servicios != nil {
		_, err = tx.Exec(`DELETE FROM gimnasio_servicios WHERE gimnasio_id = $1`, id)
		if err != nil {
			return nil, err
		}
		for _, servicio := range req.Servicios {
			_, err = tx.Exec(`
				INSERT INTO gimnasio_servicios (gimnasio_id, servicio)
				VALUES ($1, $2)
			`, id, servicio)
			if err != nil {
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return s.GetCompleto(id)
}

func (s *GimnasioService) Delete(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM gimnasio_fotos WHERE gimnasio_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM gimnasio_servicios WHERE gimnasio_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM gimnasios WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *GimnasioService) List(filters map[string]interface{}, pagina, porPagina int) ([]models.Gimnasio, int64, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM gimnasios %s", whereClause)
	err := s.db.QueryRow(countQuery, values...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	values = append(values, porPagina, offset)

	query := fmt.Sprintf(`
		SELECT id, nombre, direccion, barrio, ciudad, pais, telefono, horario, descripcion, imagen, precio, moneda, activo, lat, lng, creado_en
		FROM gimnasios %s ORDER BY creado_en DESC LIMIT $%d OFFSET $%d
	`, whereClause, i, i+1)

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var gimnasios []models.Gimnasio
	for rows.Next() {
		var g models.Gimnasio
		err := rows.Scan(&g.ID, &g.Nombre, &g.Direccion, &g.Barrio, &g.Ciudad, &g.Pais, &g.Telefono, &g.Horario, &g.Descripcion, &g.Imagen, &g.Precio, &g.Moneda, &g.Activo, &g.Lat, &g.Lng, &g.CreadoEn)
		if err != nil {
			return nil, 0, err
		}
		gimnasios = append(gimnasios, g)
	}

	return gimnasios, total, nil
}

func (s *GimnasioService) GetClasesByGimnasio(gimnasioID int) ([]models.ClaseDisponible, error) {
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