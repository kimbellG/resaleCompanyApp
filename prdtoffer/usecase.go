package prdtoffer

import (
	"context"
)

type Offer struct {
	Id           int
	ProductName  string `json:"product"`
	ProviderName string `json:"provider"`
	Cost         float32
}

type IdController interface {
	GetIDByName(ctx context.Context, name string) (int, error)
	GetNameById(id int) (string, error)
}

type UseCase interface {
	Add(ctx context.Context, pr *Offer) error
	GetOfferForProduct(ctx context.Context, productName string) ([]Offer, error)
	GetOffersOfProvider(ctx context.Context, providerName string) ([]Offer, error)
	UpdateCost(ctx context.Context, productName, providerName string, cost float32) error
	Delete(ctx context.Context, productName, providerName string) error
}
