package product

import (
	"context"
)

type Product struct {
	Id          int
	Name        string
	Description string
}

type UseCase interface {
	Add(ctx context.Context, product *Product) error
	Gets(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, id int, fields map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	Filter(ctx context.Context, key, value string) ([]Product, error)
}
