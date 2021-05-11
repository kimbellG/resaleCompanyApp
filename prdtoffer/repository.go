package prdtoffer

import (
	"context"
	"cw/models"
)

type Repository interface {
	Add(ctx context.Context, offer *models.Offer) error
	GetOfferForProduct(ctx context.Context, productId int) ([]models.Offer, error)
	GetOfferOfProvider(ctx context.Context, providerId int) ([]models.Offer, error)
	UpdateCost(ctx context.Context, providerId, productId int, cost float32) error
	Delete(ctx context.Context, providerId, productId int) error
	GetById(id int) (*models.Offer, error)
}
