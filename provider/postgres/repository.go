package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type ProviderRepository struct {
	db *sql.DB
}

func NewProviderRepository(lib_db *sql.DB) *ProviderRepository {
	err := dbutil.CreateTable(lib_db,
		`CREATE TABLE IF NOT EXISTS Provider (
	 		vendor_code		 SERIAL PRIMARY KEY UNIQUE,
	 		name			 VARCHAR(200) NOT NULL UNIQUE,
	 		unp				 VARCHAR(10) NOT NULL CHECK(char_length(unp) = 9),
	 		terms_of_payment VARCHAR(100),
			address			 VARCHAR(200) NOT NULL,
			phone_number	 CHAR(14) CHECK(char_length(phone_number) = 13),
			email			 VARCHAR(100),
			web_site		 VARCHAR(100)
	);`)

	if err != nil {
		panic(err)
	}

	return &ProviderRepository{
		db: lib_db,
	}
}

func (pr *ProviderRepository) AddProvider(ctx context.Context, mp *models.Provider) error {
	stmt, err := pr.db.Prepare(
		`INSERT INTO Provider (name, unp, terms_of_payment, address, phone_number, email, web_site) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(mp.Name, mp.UNP, mp.TermsOfPayment, mp.Address, mp.PhoneNumber, mp.Email, mp.WebSite); err != nil {
		return err
	}

	return nil
}

func (pr *ProviderRepository) UpdateProvider(ctx context.Context, code int, fields map[string]interface{}) error {
	for key, value := range fields {
		stmt, err := pr.db.Prepare(fmt.Sprintf("UPDATE provider SET %v=$1 WHERE vendor_code = $2", key))
		if err != nil {
			return fmt.Errorf("incorrect field in provider table")
		}

		if _, err := stmt.Exec(value, code); err != nil {
			return fmt.Errorf("incorrect value this %v key", key)
		}
	}
	return nil

}

func (pr *ProviderRepository) DeleteProvider(ctx context.Context, code int) error {
	stmt, err := pr.db.Prepare("DELETE FROM Provider WHERE vendor_code = $1")
	if err != nil {
		return fmt.Errorf("incorrect stmt: %v", err)
	}

	if _, err = stmt.Exec(code); err != nil {
		return fmt.Errorf("invalid code %v: %v", code, err)
	}

	return nil
}

func (pr *ProviderRepository) GetProviders(ctx context.Context) ([]models.Provider, error) {
	stmt, err := pr.db.Prepare("SELECT * FROM Provider")
	if err != nil {
		return nil, fmt.Errorf("incorrect stmt: %v", err)
	}

	dbRlt, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("db query is failed: %v", err)
	}
	defer dbRlt.Close()

	result := new([]models.Provider)
	for dbRlt.Next() {
		tmp := models.Provider{}
		if err := dbRlt.Scan(&tmp.VendorCode, &tmp.Name, &tmp.UNP, &tmp.TermsOfPayment, &tmp.Address, &tmp.PhoneNumber, &tmp.Email, &tmp.WebSite); err != nil {
			return nil, fmt.Errorf("scan is failed: %v", err)
		}
		*result = append(*result, tmp)
	}

	return *result, nil
}

func (pr *ProviderRepository) DeleteAll() error {
	stmt, err := pr.db.Prepare("DELETE FROM Provider")
	if err != nil {
		return fmt.Errorf("incorrect stmt: %v", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	return nil

}
