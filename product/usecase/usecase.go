package usecase

import (
	"context"
	"cw/models"
	"cw/product"
	"fmt"
	"reflect"
	"strings"
)

type ProductUseCase struct {
	rep product.Repository
}

func NewProductUseCase(rep_pr product.Repository) *ProductUseCase {
	return &ProductUseCase{
		rep: rep_pr,
	}
}

func (p *ProductUseCase) Add(ctx context.Context, pr *product.Product) error {
	modProvider := clToModels(pr)

	return p.rep.Add(ctx, modProvider)
}

func clToModels(pr *product.Product) *models.Product {
	return &models.Product{
		Id:          pr.Id,
		Name:        pr.Name,
		Description: pr.Description,
	}
}

func (p *ProductUseCase) Gets(ctx context.Context) ([]product.Product, error) {
	modClient, err := p.rep.Gets(ctx)
	if err != nil {
		return nil, err
	}

	return arrModToPr(modClient), nil
}

func arrModToPr(mods []models.Product) []product.Product {
	result := new([]product.Product)
	for _, mod := range mods {
		*result = append(*result, *modToCl(&mod))
	}

	return *result
}

func modToCl(mod *models.Product) *product.Product {
	return &product.Product{
		Id:          mod.Id,
		Name:        mod.Name,
		Description: mod.Description,
	}
}

func (p *ProductUseCase) Update(ctx context.Context, code int, fields map[string]interface{}) error {
	return p.rep.Update(ctx, code, fields)
}

func (p *ProductUseCase) Delete(ctx context.Context, code int) error {
	return p.rep.Delete(ctx, code)
}

func (u *ProductUseCase) Filter(ctx context.Context, key, value string) ([]product.Product, error) {
	if isNotKeyOfStruct(key) {
		return nil, fmt.Errorf("incorrect name of field: %v", key)
	}

	result, err := u.rep.Filter(ctx, key, value)
	if err != nil {
		return nil, err
	}

	return arrModToPr(result), nil
}

func isNotKeyOfStruct(key string) bool {
	pr := &product.Product{}
	s := reflect.ValueOf(pr).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if strings.ToLower(typeOfT.Field(i).Name) == key {
			return false
		}
	}

	return true
}
