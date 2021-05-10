package order

import (
	"context"
	"cw/client"
	"cw/models"
	"cw/prdtoffer"
	"time"
)

type OrderOutput struct {
	Id        int
	Offers    []prdtoffer.Offer
	Client    client.Client
	Manager   string
	OrderDate time.Time
	Quantity  int
	Status    string
}

type OfferController interface {
	GetById(id int) (*prdtoffer.Offer, error)
}

type ClientController interface {
	GetById(id int) (*client.Client, error)
}

type ManagerController interface {
	GetNameByLogin(login string) (string, error)
}

type UseCase interface {
	Add(ctx context.Context, order *models.Order) error
	Gets(ctx context.Context) ([]OrderOutput, error)
	GetInInterval(ctx context.Context, start, end time.Time) ([]OrderOutput, error)
	UpdateStatus(ctx context.Context, id int, newStatus string) error
	Filter(ctx context.Context, key string, value interface{}) ([]OrderOutput, error)
}
