package client

import (
	"context"
)

type Client struct {
	Id           int
	Name         string
	FIO          string
	Address      string
	Phone_Number string
	Email        string
}

type UseCase interface {
	AddClient(ctx context.Context, p *Client) error
	GetClients(ctx context.Context) ([]Client, error)
	UpdateClient(ctx context.Context, code int, fields map[string]interface{}) error
	DeleteClient(ctx context.Context, code int) error
	FilterClient(ctx context.Context, key, value string) ([]Client, error)
}
