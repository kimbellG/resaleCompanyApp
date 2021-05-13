package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(lib_db *sql.DB) *ProductRepository {
	err := dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS Product (
	 		id		 SERIAL PRIMARY KEY UNIQUE,
	 		Name			 VARCHAR(200) NOT NULL UNIQUE,
			Description 	 VARCHAR(1000)
	);`)

	if err != nil {
		panic(err)
	}

	return &ProductRepository{
		db: lib_db,
	}
}

func (p *ProductRepository) Add(ctx context.Context, pr *models.Product) error {
	stmt, err := p.db.Prepare(
		`INSERT INTO Product (Name, Description) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("prepare query: %v", err)
	}

	if _, err := stmt.Exec(pr.Name, pr.Description); err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

func (pr *ProductRepository) Update(ctx context.Context, code int, fields map[string]interface{}) error {
	for key, value := range fields {
		stmt, err := pr.db.Prepare(fmt.Sprintf("UPDATE Product SET %v=$1 WHERE id = $2", key))
		if err != nil {
			return fmt.Errorf("incorrect field in provider table: %v", err)
		}

		if _, err := stmt.Exec(value, code); err != nil {
			return fmt.Errorf("incorrect value this %v key", key)
		}
	}
	return nil
}

func (pr *ProductRepository) Delete(ctx context.Context, code int) error {
	stmt, err := pr.db.Prepare("DELETE FROM Product WHERE id = $1")
	if err != nil {
		return fmt.Errorf("incorrect stmt: %v", err)
	}

	if _, err = stmt.Exec(code); err != nil {
		return fmt.Errorf("invalid code %v: %v", code, err)
	}

	return nil
}

func (pr *ProductRepository) Gets(ctx context.Context) ([]models.Product, error) {
	stmt, err := pr.db.Prepare("SELECT * FROM Product")
	if err != nil {
		return nil, fmt.Errorf("incorrect stmt: %v", err)
	}

	dbRlt, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("db query is failed: %v", err)
	}
	defer dbRlt.Close()

	result := new([]models.Product)

	*result, err = scanClient(dbRlt, *result)
	if err != nil {
		return nil, err
	}

	return *result, nil
}

func (pr *ProductRepository) DeleteAll() error {
	stmt, err := pr.db.Prepare("DELETE FROM Product")
	if err != nil {
		return fmt.Errorf("incorrect stmt: %v", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Filter(ctx context.Context, key, value string) ([]models.Product, error) {
	stmt, err := r.db.Prepare(fmt.Sprintf("SELECT * FROM Product WHERE %v LIKE $1", key))
	if err != nil {
		return nil, fmt.Errorf("incorrect stmt to db: %v", err)
	}

	query, err := stmt.Query(fmt.Sprintf("%%%v%%", value))
	if err != nil {
		return nil, fmt.Errorf("stmt-query error: %v", err)
	}

	result := make([]models.Product, 0)

	result, err = scanClient(query, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func scanClient(rows *sql.Rows, result []models.Product) ([]models.Product, error) {
	for rows.Next() {
		tmp := models.Product{}
		if err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Description); err != nil {
			return nil, fmt.Errorf("scan is failed: %v", err)
		}

		if tmp.Id == 0 {
			continue
		}

		result = append(result, tmp)
	}

	return result, nil
}

func (r *ProductRepository) GetIDByName(ctx context.Context, name string) (int, error) {
	stmt, err := r.db.Prepare("SELECT id FROM Product WHERE name = $1")
	if err != nil {
		return -1, fmt.Errorf("prepare stmt: %v", err)
	}

	query, err := stmt.Query(name)
	if err != nil {
		return -1, fmt.Errorf("exec stmt: %v", err)
	}

	var result int
	for query.Next() {
		if err := query.Scan(&result); err != nil {
			return -1, fmt.Errorf("scan: %v", err)
		}
	}

	if result == 0 {
		return -1, fmt.Errorf("result is empty")
	}

	return result, nil
}

func (r *ProductRepository) GetNameById(id int) (string, error) {
	stmt, err := r.db.Prepare("SELECT name FROM Product WHERE id = $1")
	if err != nil {
		return "", fmt.Errorf("prepare stmt: %v", err)
	}

	query, err := stmt.Query(id)
	if err != nil {
		return "", fmt.Errorf("query stmt: %v", err)
	}

	var name string
	for query.Next() {
		if err := query.Scan(&name); err != nil {
			return "", fmt.Errorf("scan: %v", err)
		}
	}

	if name == "" {
		return "", fmt.Errorf("name not found")
	}

	return name, nil
}
