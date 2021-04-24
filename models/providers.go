package models

type ProviderJSON struct {
	Reg  Region
	Prov Provider
}

type Provider struct {
	VendorCode     int
	Name           string
	TermsOfPayment string
	RegionCode     string
}

type Region struct {
	Country     string
	City        string
	Address     string
	PhoneNumber string
	Email       string
}
