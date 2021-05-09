package client

import (
	"context"
	"cw/models"
)

type Repository interface {
	AddClient(ctx context.Context, pr *models.Client) error
	GetClients(ctx context.Context) ([]models.Client, error)
	UpdateClient(ctx context.Context, code int, fields map[string]interface{}) error
	DeleteClient(ctx context.Context, code int) error
	FilterClient(ctx context.Context, key, value string) ([]models.Client, error)
}
