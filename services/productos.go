package services

import (
	"database/sql"
	"fmt"
	"strings"

	"api_mid_the_housefit/models"
)

type ProductoService struct {
	db *sql.DB
}

func NewProductoService() *ProductoService {
	return &ProductoService{db: GetDB()}
}

func (s *ProductoService) Create(req models.ProductoRequest) (*models.Producto, error) {
	var producto models.Producto
	err := s.db.QueryRow(`
		INSERT INTO productos (nombre, categoria, precio, stock, descripcion, guia_nutricional, imagen)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, nombre, categoria, precio, stock, descripcion, guia_nutricional, imagen, creado_en
	`, req.Nombre, req.Categoria, req.Precio, req.Stock, req.Descripcion, req.GuiaNutricional, req.Imagen).Scan(
		&producto.ID, &producto.Nombre, &producto.Categoria, &producto.Precio, &producto.Stock,
		&producto.Descripcion, &producto.GuiaNutricional, &producto.Imagen, &producto.CreadoEn,
	)
	return &producto, err
}

func (s *ProductoService) GetByID(id int) (*models.Producto, error) {
	var producto models.Producto
	err := s.db.QueryRow(`
		SELECT id, nombre, categoria, precio, stock, descripcion, guia_nutricional, imagen, creado_en
		FROM productos WHERE id = $1
	`, id).Scan(
		&producto.ID, &producto.Nombre, &producto.Categoria, &producto.Precio, &producto.Stock,
		&producto.Descripcion, &producto.GuiaNutricional, &producto.Imagen, &producto.CreadoEn,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &producto, err
}

func (s *ProductoService) Update(id int, req models.ProductoUpdateRequest) (*models.Producto, error) {
	setClauses := []string{}
	values := []interface{}{}
	i := 1

	if req.Nombre != nil {
		setClauses = append(setClauses, fmt.Sprintf("nombre = $%d", i))
		values = append(values, *req.Nombre)
		i++
	}
	if req.Categoria != nil {
		setClauses = append(setClauses, fmt.Sprintf("categoria = $%d", i))
		values = append(values, *req.Categoria)
		i++
	}
	if req.Precio != nil {
		setClauses = append(setClauses, fmt.Sprintf("precio = $%d", i))
		values = append(values, *req.Precio)
		i++
	}
	if req.Stock != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock = $%d", i))
		values = append(values, *req.Stock)
		i++
	}
	if req.Descripcion != nil {
		setClauses = append(setClauses, fmt.Sprintf("descripcion = $%d", i))
		values = append(values, *req.Descripcion)
		i++
	}
	if req.GuiaNutricional != nil {
		setClauses = append(setClauses, fmt.Sprintf("guia_nutricional = $%d", i))
		values = append(values, *req.GuiaNutricional)
		i++
	}
	if req.Imagen != nil {
		setClauses = append(setClauses, fmt.Sprintf("imagen = $%d", i))
		values = append(values, *req.Imagen)
		i++
	}

	if len(setClauses) == 0 {
		return s.GetByID(id)
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE productos SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)
	_, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *ProductoService) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM productos WHERE id = $1`, id)
	return err
}

func (s *ProductoService) List(filters map[string]interface{}, pagina, porPagina int) ([]models.Producto, int64, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM productos %s", whereClause)
	err := s.db.QueryRow(countQuery, values...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina
	values = append(values, porPagina, offset)

	query := fmt.Sprintf(`
		SELECT id, nombre, categoria, precio, stock, descripcion, guia_nutricional, imagen, creado_en
		FROM productos %s ORDER BY creado_en DESC LIMIT $%d OFFSET $%d
	`, whereClause, i, i+1)

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		err := rows.Scan(&p.ID, &p.Nombre, &p.Categoria, &p.Precio, &p.Stock, &p.Descripcion, &p.GuiaNutricional, &p.Imagen, &p.CreadoEn)
		if err != nil {
			return nil, 0, err
		}
		productos = append(productos, p)
	}

	return productos, total, nil
}
