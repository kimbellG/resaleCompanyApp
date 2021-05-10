package order

import (
	"context"
	"cw/models"
)

type Repository interface {
	Add(ctx context.Context, order *models.Order) error
	Gets(ctx context.Context) ([]models.Order, error)
	GetInInterval(ctx context.Context, start, end string) ([]models.Order, error)
	UpdateStatus(ctx context.Context, id int, newStatus string) error
	Filter(ctx context.Context, key string, value interface{}) ([]models.Order, error)
}
