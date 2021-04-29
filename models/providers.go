package models

import (
	"fmt"

	"github.com/pkg/errors"
)

func (db *DBController) InsertProviderInDB(pr *Provider) error {
	stmt, err := db.Prepare(
		`INSERT INTO Provider (name, unp, terms_of_payment, address, phone_number, email, web_site) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(pr.Name, pr.UNP, pr.TermsOfPayment, pr.Address, pr.PhoneNumber, pr.Email, pr.WebSite); err != nil {
		return err
	}

	return nil
}

func (db *DBController) EditFieldsInProviderTable(vendorCode int, fields map[string]string) error {
	for key, value := range fields {
		stmt, err := db.Prepare(fmt.Sprintf("UPDATE provider SET %v=$1 WHERE vendor_code = $2", key))
		if err != nil {
			return errors.Errorf("Incorrect field in provider table")
		}

		if _, err := stmt.Exec(value, vendorCode); err != nil {
			return errors.Errorf("Incorrect value this %v key", key)
		}
	}
	return nil
}
