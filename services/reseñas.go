package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_housefit/models"
)

type ResenaService struct {
	db *sql.DB
}

func NewResenaService() *ResenaService {
	return &ResenaService{db: GetDB()}
}

func (s *ResenaService) Create(req models.ResenaRequest) (*models.ResenaResponse, error) {
	var resena models.Resena
	err := s.db.QueryRow(`
		INSERT INTO resenas (gimnasio_id, usuario_correo, usuario_nombre, titulo, cuerpo, imagen, puntuacion)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, gimnasio_id, usuario_correo, usuario_nombre, titulo, cuerpo, imagen, puntuacion, fecha
	`, req.GimnasioID, req.UsuarioCorreo, req.UsuarioNombre, req.Titulo, req.Cuerpo, req.Imagen, req.Puntuacion).Scan(
		&resena.ID, &resena.GimnasioID, &resena.UsuarioCorreo, &resena.UsuarioNombre,
		&resena.Titulo, &resena.Cuerpo, &resena.Imagen, &resena.Puntuacion, &resena.Fecha,
	)
	if err != nil {
		return nil, err
	}

	gimnasioService := NewGimnasioService()
	gimnasio, _ := gimnasioService.GetCompleto(req.GimnasioID)

	var g *models.Gimnasio
	if gimnasio != nil {
		g = &gimnasio.Gimnasio
	}

	return &models.ResenaResponse{
		Resena:   resena,
		Gimnasio: g,
	}, nil
}

func (s *ResenaService) GetByID(id int) (*models.ResenaResponse, error) {
	var resena models.Resena
	err := s.db.QueryRow(`
		SELECT id, gimnasio_id, usuario_correo, usuario_nombre, titulo, cuerpo, imagen, puntuacion, fecha
		FROM resenas WHERE id = $1
	`, id).Scan(&resena.ID, &resena.GimnasioID, &resena.UsuarioCorreo, &resena.UsuarioNombre, &resena.Titulo, &resena.Cuerpo, &resena.Imagen, &resena.Puntuacion, &resena.Fecha)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	gimnasioService := NewGimnasioService()
	gimnasio, _ := gimnasioService.GetCompleto(resena.GimnasioID)

	var g *models.Gimnasio
	if gimnasio != nil {
		g = &gimnasio.Gimnasio
	}

	return &models.ResenaResponse{
		Resena:   resena,
		Gimnasio: g,
	}, nil
}

func (s *ResenaService) Update(id int, req models.ResenaUpdateRequest) (*models.ResenaResponse, error) {
	setClauses := []string{}
	values := []interface{}{}
	i := 1

	if req.Titulo != nil {
		setClauses = append(setClauses, fmt.Sprintf("titulo = $%d", i))
		values = append(values, *req.Titulo)
		i++
	}
	if req.Cuerpo != nil {
		setClauses = append(setClauses, fmt.Sprintf("cuerpo = $%d", i))
		values = append(values, *req.Cuerpo)
		i++
	}
	if req.Imagen != nil {
		setClauses = append(setClauses, fmt.Sprintf("imagen = $%d", i))
		values = append(values, *req.Imagen)
		i++
	}
	if req.Puntuacion != nil {
		setClauses = append(setClauses, fmt.Sprintf("puntuacion = $%d", i))
		values = append(values, *req.Puntuacion)
		i++
	}

	if len(setClauses) == 0 {
		return s.GetByID(id)
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE resenas SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)
	_, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *ResenaService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM resenas WHERE id = $1`, id)
	return err
}

func (s *ResenaService) ListByGimnasio(gimnasioID int, pagina, porPagina int) ([]models.ResenaResponse, int64, error) {
	var total int64
	err := s.db.QueryRow(`SELECT COUNT(*) FROM resenas WHERE gimnasio_id = $1`, gimnasioID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	rows, err := s.db.Query(`
		SELECT id, gimnasio_id, usuario_correo, usuario_nombre, titulo, cuerpo, imagen, puntuacion, fecha
		FROM resenas WHERE gimnasio_id = $1 ORDER BY fecha DESC LIMIT $2 OFFSET $3
	`, gimnasioID, porPagina, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var resenas []models.ResenaResponse
	for rows.Next() {
		var r models.Resena
		err := rows.Scan(&r.ID, &r.GimnasioID, &r.UsuarioCorreo, &r.UsuarioNombre, &r.Titulo, &r.Cuerpo, &r.Imagen, &r.Puntuacion, &r.Fecha)
		if err != nil {
			return nil, 0, err
		}

		gimnasioService := NewGimnasioService()
		gimnasio, _ := gimnasioService.GetCompleto(r.GimnasioID)

		var g *models.Gimnasio
		if gimnasio != nil {
			g = &gimnasio.Gimnasio
		}

		resenas = append(resenas, models.ResenaResponse{
			Resena:   r,
			Gimnasio: g,
		})
	}

	return resenas, total, nil
}

func (s *ResenaService) GetPromedioGimnasio(gimnasioID int) (float64, int, error) {
	var promedio sql.NullFloat64
	var total int
	err := s.db.QueryRow(`
		SELECT AVG(puntuacion)::numeric(3,2), COUNT(*) FROM resenas WHERE gimnasio_id = $1
	`, gimnasioID).Scan(&promedio, &total)
	if err != nil {
		return 0, 0, err
	}
	if promedio.Valid {
		return float64(promedio.Float64), total, nil
	}
	return 0, total, nil
}
