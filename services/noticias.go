package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_housefit/models"
)

type NoticiaService struct {
	db *sql.DB
}

func NewNoticiaService() *NoticiaService {
	return &NoticiaService{db: GetDB()}
}

func (s *NoticiaService) List(filters map[string]interface{}, pagina, porPagina int, soloActivas bool) ([]models.Noticia, int64, error) {
	if s.db == nil {
		return []models.Noticia{}, 0, nil
	}

	whereClauses := []string{}
	args := []interface{}{}
	idx := 1

	if soloActivas {
		whereClauses = append(whereClauses, fmt.Sprintf("activa = $%d", idx))
		args = append(args, true)
		idx++
	}
	if v, ok := filters["tipo"]; ok && v != nil && fmt.Sprint(v) != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("tipo = $%d", idx))
		args = append(args, v)
		idx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM noticias" + whereSQL
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	query := "SELECT id, titulo, slug, resumen, contenido, imagen, tipo, activa, creado_en FROM noticias" + whereSQL + " ORDER BY creado_en DESC LIMIT $" + fmt.Sprint(idx) + " OFFSET $" + fmt.Sprint(idx+1)
	args = append(args, porPagina, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var noticias []models.Noticia
	for rows.Next() {
		var n models.Noticia
		if err := rows.Scan(&n.ID, &n.Titulo, &n.Slug, &n.Resumen, &n.Contenido, &n.Imagen, &n.Tipo, &n.Activa, &n.CreadoEn); err != nil {
			return nil, 0, err
		}
		noticias = append(noticias, n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return noticias, total, nil
}

func (s *NoticiaService) GetByID(id int) (*models.Noticia, error) {
	if s.db == nil {
		return nil, nil
	}
	var n models.Noticia
	err := s.db.QueryRow(`SELECT id, titulo, slug, resumen, contenido, imagen, tipo, activa, creado_en FROM noticias WHERE id = $1`, id).Scan(
		&n.ID, &n.Titulo, &n.Slug, &n.Resumen, &n.Contenido, &n.Imagen, &n.Tipo, &n.Activa, &n.CreadoEn,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *NoticiaService) Create(req models.NoticiaRequest) (*models.Noticia, error) {
	if s.db == nil {
		return nil, nil
	}
	var n models.Noticia
	err := s.db.QueryRow(`
		INSERT INTO noticias (titulo, slug, resumen, contenido, imagen, tipo, activa)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, titulo, slug, resumen, contenido, imagen, tipo, activa, creado_en
	`, req.Titulo, req.Slug, req.Resumen, req.Contenido, req.Imagen, req.Tipo, req.Activa).Scan(
		&n.ID, &n.Titulo, &n.Slug, &n.Resumen, &n.Contenido, &n.Imagen, &n.Tipo, &n.Activa, &n.CreadoEn,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *NoticiaService) Update(id int, req models.NoticiaUpdateRequest) (*models.Noticia, error) {
	if s.db == nil {
		return nil, nil
	}
	updates := []string{}
	args := []interface{}{}
	idx := 1
	if req.Titulo != nil {
		updates = append(updates, fmt.Sprintf("titulo = $%d", idx))
		args = append(args, *req.Titulo)
		idx++
	}
	if req.Slug != nil {
		updates = append(updates, fmt.Sprintf("slug = $%d", idx))
		args = append(args, *req.Slug)
		idx++
	}
	if req.Resumen != nil {
		updates = append(updates, fmt.Sprintf("resumen = $%d", idx))
		args = append(args, *req.Resumen)
		idx++
	}
	if req.Contenido != nil {
		updates = append(updates, fmt.Sprintf("contenido = $%d", idx))
		args = append(args, *req.Contenido)
		idx++
	}
	if req.Imagen != nil {
		updates = append(updates, fmt.Sprintf("imagen = $%d", idx))
		args = append(args, *req.Imagen)
		idx++
	}
	if req.Tipo != nil {
		updates = append(updates, fmt.Sprintf("tipo = $%d", idx))
		args = append(args, *req.Tipo)
		idx++
	}
	if req.Activa != nil {
		updates = append(updates, fmt.Sprintf("activa = $%d", idx))
		args = append(args, *req.Activa)
		idx++
	}
	if len(updates) == 0 {
		return s.GetByID(id)
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE noticias SET %s WHERE id = $%d", strings.Join(updates, ", "), idx)
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *NoticiaService) Delete(id int) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.Exec(`DELETE FROM noticias WHERE id = $1`, id)
	return err
}
