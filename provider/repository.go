package provider

import (
	"context"
	"cw/models"
)

type Repository interface {
	AddProvider(ctx context.Context, pr *models.Provider) error
	GetProviders(ctx context.Context) ([]models.Provider, error)
	UpdateProvider(ctx context.Context, code int, fields map[string]interface{}) error
	DeleteProvider(ctx context.Context, code int) error
}
