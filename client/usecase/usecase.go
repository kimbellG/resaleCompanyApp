package usecase

import (
	"context"
	"cw/client"
	"cw/models"
	"fmt"
	"reflect"
	"strings"
)

type ClientUseCase struct {
	rep client.Repository
}

func NewClientUseCase(rep_pr client.Repository) *ClientUseCase {
	return &ClientUseCase{
		rep: rep_pr,
	}
}

func (p *ClientUseCase) AddClient(ctx context.Context, pr *client.Client) error {
	modProvider := prToModels(pr)

	return p.rep.AddClient(ctx, modProvider)
}

func prToModels(pr *client.Client) *models.Client {
	return &models.Client{
		Id:          pr.Id,
		Name:        pr.Name,
		FIO:         pr.FIO,
		PhoneNumber: pr.PhoneNumber,
		Address:     pr.Address,
		Email:       pr.Email,
	}
}

func (p *ClientUseCase) GetClients(ctx context.Context) ([]client.Client, error) {
	modClient, err := p.rep.GetClients(ctx)
	if err != nil {
		return nil, err
	}

	return arrModToCl(modClient), nil
}

func arrModToCl(mods []models.Client) []client.Client {
	result := new([]client.Client)
	for _, mod := range mods {
		*result = append(*result, *modToPr(&mod))
	}

	return *result
}

func modToPr(mod *models.Client) *client.Client {
	return &client.Client{
		Id:          mod.Id,
		Name:        mod.Name,
		FIO:         mod.FIO,
		PhoneNumber: mod.PhoneNumber,
		Address:     mod.Address,
		Email:       mod.Email,
	}
}

func (p *ClientUseCase) UpdateClient(ctx context.Context, code int, fields map[string]interface{}) error {
	return p.rep.UpdateClient(ctx, code, fields)
}

func (p *ClientUseCase) DeleteClient(ctx context.Context, code int) error {
	return p.rep.DeleteClient(ctx, code)
}

func (u *ClientUseCase) FilterClient(ctx context.Context, key, value string) ([]client.Client, error) {
	if isNotKeyOfStruct(key) {
		return nil, fmt.Errorf("incorrect name of field: %v", key)
	}

	result, err := u.rep.FilterClient(ctx, key, value)
	if err != nil {
		return nil, err
	}

	return arrModToCl(result), nil
}

func isNotKeyOfStruct(key string) bool {
	pr := &client.Client{}
	s := reflect.ValueOf(pr).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if strings.ToLower(typeOfT.Field(i).Name) == key {
			return false
		}
	}

	return true
}
