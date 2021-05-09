package product

import (
	"context"
	"cw/models"
)

type Repository interface {
	Add(ctx context.Context, pr *models.Product) error
	Gets(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, id int, fields map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	Filter(ctx context.Context, key string, value string) ([]models.Product, error)
}
