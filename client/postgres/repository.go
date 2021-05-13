package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(lib_db *sql.DB) *ClientRepository {
	err := dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS Client (
	 		id		 SERIAL PRIMARY KEY UNIQUE,
	 		Name			 VARCHAR(200) NOT NULL,
	 		FIO				 VARCHAR(200) NOT NULL,
			address			 VARCHAR(200) NOT NULL,
			email			 VARCHAR(100),
			phone_number	 	 CHAR(14) CHECK(char_length(phone_number) = 13)
	);`)

	if err != nil {
		panic(err)
	}

	return &ClientRepository{
		db: lib_db,
	}
}

func (pr *ClientRepository) AddClient(ctx context.Context, mp *models.Client) error {
	stmt, err := pr.db.Prepare(
		`INSERT INTO Client (Name, FIO, address, phone_number, email) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(mp.Name, mp.FIO, mp.Address, mp.PhoneNumber, mp.Email); err != nil {
		return err
	}

	return nil
}

func (pr *ClientRepository) UpdateClient(ctx context.Context, code int, fields map[string]interface{}) error {
	for key, value := range fields {
		stmt, err := pr.db.Prepare(fmt.Sprintf("UPDATE Client SET %v=$1 WHERE id = $2", key))
		if err != nil {
			return fmt.Errorf("incorrect field in provider table: %v", err)
		}

		if _, err := stmt.Exec(value, code); err != nil {
			return fmt.Errorf("incorrect value this %v key", key)
		}
	}
	return nil

}

func (pr *ClientRepository) DeleteClient(ctx context.Context, code int) error {
	stmt, err := pr.db.Prepare("DELETE FROM Client WHERE id = $1")
	if err != nil {
		return fmt.Errorf("incorrect stmt: %v", err)
	}

	if _, err = stmt.Exec(code); err != nil {
		return fmt.Errorf("invalid code %v: %v", code, err)
	}

	return nil
}

func (pr *ClientRepository) GetClients(ctx context.Context) ([]models.Client, error) {
	stmt, err := pr.db.Prepare("SELECT * FROM Client")
	if err != nil {
		return nil, fmt.Errorf("incorrect stmt: %v", err)
	}

	dbRlt, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("db query is failed: %v", err)
	}
	defer dbRlt.Close()

	result := new([]models.Client)

	*result, err = scanClient(dbRlt, *result)
	if err != nil {
		return nil, err
	}

	return *result, nil
}

func (pr *ClientRepository) DeleteAll() error {
	stmt, err := pr.db.Prepare("DELETE FROM Client")
	if err != nil {
		return fmt.Errorf("incorrect stmt: %v", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepository) FilterClient(ctx context.Context, key, value string) ([]models.Client, error) {
	stmt, err := r.db.Prepare(fmt.Sprintf("SELECT * FROM Client WHERE %v LIKE $1", key))
	if err != nil {
		return nil, fmt.Errorf("incorrect stmt to db: %v", err)
	}

	query, err := stmt.Query(fmt.Sprintf("%%%v%%", value))
	if err != nil {
		return nil, fmt.Errorf("stmt-query error: %v", err)
	}

	result := make([]models.Client, 0)

	result, err = scanClient(query, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func scanClient(rows *sql.Rows, result []models.Client) ([]models.Client, error) {
	for rows.Next() {
		tmp := models.Client{}
		if err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.FIO, &tmp.Address, &tmp.PhoneNumber, &tmp.Email); err != nil {
			return nil, fmt.Errorf("scan is failed: %v", err)
		}

		if tmp.Id == 0 {
			continue
		}

		result = append(result, tmp)
	}

	return result, nil
}
