package models

type Provider struct {
	VendorCode     int    `json:"vendor_code"`
	Name           string `json:"name"`
	UNP            string `json:"unp"`
	TermsOfPayment string `json:"terms_of_payment"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	WebSite        string `json:"web_site"`
}
