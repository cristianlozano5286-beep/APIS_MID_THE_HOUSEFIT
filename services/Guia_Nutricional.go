package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_house_fit/models"
)

type GuiaNutricionalService struct {
	db *sql.DB
}

func NewGuiaNutricionalService() *GuiaNutricionalService {
	return &GuiaNutricionalService{db: GetDB()}
}

func (s *GuiaNutricionalService) Create(req models.GuiaNutricionalRequest) (*models.GuiaNutricional, error) {
	var guia models.GuiaNutricional
	err := s.db.QueryRow(`
		INSERT INTO guias_nutricionales (categoria_imc, titulo, recomendaciones, ejemplo_comidas)
		VALUES ($1, $2, $3, $4)
		RETURNING id, categoria_imc, titulo, recomendaciones, ejemplo_comidas, creado_en
	`, req.CategoriaIMC, req.Titulo, req.Recomendaciones, req.EjemploComidas).Scan(
		&guia.ID, &guia.CategoriaIMC, &guia.Titulo, &guia.Recomendaciones, &guia.EjemploComidas, &guia.CreadoEn,
	)
	return &guia, err
}

func (s *GuiaNutricionalService) GetByID(id int) (*models.GuiaNutricional, error) {
	var guia models.GuiaNutricional
	err := s.db.QueryRow(`
		SELECT id, categoria_imc, titulo, recomendaciones, ejemplo_comidas, creado_en
		FROM guias_nutricionales WHERE id = $1
	`, id).Scan(&guia.ID, &guia.CategoriaIMC, &guia.Titulo, &guia.Recomendaciones, &guia.EjemploComidas, &guia.CreadoEn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &guia, err
}

func (s *GuiaNutricionalService) GetByCategoria(categoria string) (*models.GuiaNutricional, error) {
	var guia models.GuiaNutricional
	err := s.db.QueryRow(`
		SELECT id, categoria_imc, titulo, recomendaciones, ejemplo_comidas, creado_en
		FROM guias_nutricionales WHERE categoria_imc = $1
	`, categoria).Scan(&guia.ID, &guia.CategoriaIMC, &guia.Titulo, &guia.Recomendaciones, &guia.EjemploComidas, &guia.CreadoEn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &guia, err
}

func (s *GuiaNutricionalService) Update(id int, req models.GuiaNutricionalUpdateRequest) (*models.GuiaNutricional, error) {
	setClauses := []string{}
	values := []interface{}{}
	i := 1

	if req.CategoriaIMC != nil {
		setClauses = append(setClauses, fmt.Sprintf("categoria_imc = $%d", i))
		values = append(values, *req.CategoriaIMC)
		i++
	}
	if req.Titulo != nil {
		setClauses = append(setClauses, fmt.Sprintf("titulo = $%d", i))
		values = append(values, *req.Titulo)
		i++
	}
	if req.Recomendaciones != nil {
		setClauses = append(setClauses, fmt.Sprintf("recomendaciones = $%d", i))
		values = append(values, *req.Recomendaciones)
		i++
	}
	if req.EjemploComidas != nil {
		setClauses = append(setClauses, fmt.Sprintf("ejemplo_comidas = $%d", i))
		values = append(values, *req.EjemploComidas)
		i++
	}

	if len(setClauses) == 0 {
		return s.GetByID(id)
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE guias_nutricionales SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)
	_, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *GuiaNutricionalService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM guias_nutricionales WHERE id = $1`, id)
	return err
}

func (s *GuiaNutricionalService) List(pagina, porPagina int) ([]models.GuiaNutricional, int64, error) {
	var total int64
	err := s.db.QueryRow(`SELECT COUNT(*) FROM guias_nutricionales`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	rows, err := s.db.Query(`
		SELECT id, categoria_imc, titulo, recomendaciones, ejemplo_comidas, creado_en
		FROM guias_nutricionales ORDER BY categoria_imc, creado_en DESC LIMIT $1 OFFSET $2
	`, porPagina, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var guias []models.GuiaNutricional
	for rows.Next() {
		var g models.GuiaNutricional
		err := rows.Scan(&g.ID, &g.CategoriaIMC, &g.Titulo, &g.Recomendaciones, &g.EjemploComidas, &g.CreadoEn)
		if err != nil {
			return nil, 0, err
		}
		guias = append(guias, g)
	}

	return guias, total, nil
}