package provider

import (
	"context"
)

type Provider struct {
	VendorCode     int
	Name           string
	UNP            string
	TermsOfPayment string
	Address        string
	PhoneNumber    string
	Email          string
	WebSite        string
}

type UseCase interface {
	AddProvider(ctx context.Context, p *Provider) error
	GetProviders(ctx context.Context) ([]Provider, error)
	UpdateProvider(ctx context.Context, code int, fields map[string]interface{}) error
	DeleteProvider(ctx context.Context, code int) error
}
