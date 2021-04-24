package models

type Provider struct {
	VendorCode     int
	Name           string
	UNP            string `json:"unp"`
	TermsOfPayment string `json:"terms_of_payment"`
	Address        string
	PhoneNumber    string `json:"phone_number"`
	Email          string
	WebSite        string `json:"web_site"`
}

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
