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
