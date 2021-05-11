package usecase

import (
	"context"
	"cw/models"
	"cw/prdtoffer"
	"fmt"
)

type ProductOfferUseCase struct {
	rep  prdtoffer.Repository
	prov prdtoffer.IdController
	prod prdtoffer.IdController
}

func NewProductOfferUseCase(repo prdtoffer.Repository, provider prdtoffer.IdController, product prdtoffer.IdController) *ProductOfferUseCase {
	return &ProductOfferUseCase{
		rep:  repo,
		prov: provider,
		prod: product,
	}
}

func (o *ProductOfferUseCase) Add(ctx context.Context, pr *prdtoffer.Offer) error {
	newOffer := &models.Offer{}

	if err := o.feelProviderId(ctx, pr.ProviderName, newOffer); err != nil {
		return err
	}
	if err := o.feelProductId(ctx, pr.ProductName, newOffer); err != nil {
		return err
	}
	newOffer.Cost = pr.Cost

	if err := o.rep.Add(ctx, newOffer); err != nil {
		return err
	}

	return nil
}

func (o *ProductOfferUseCase) feelProviderId(ctx context.Context, name string, offer *models.Offer) error {
	providerId, err := o.prov.GetIDByName(ctx, name)
	if err != nil {
		return fmt.Errorf("incorrect name of provider: %v", err)
	}

	offer.ProviderId = providerId
	return nil
}

func (o *ProductOfferUseCase) feelProductId(ctx context.Context, name string, offer *models.Offer) error {
	productId, err := o.prod.GetIDByName(ctx, name)
	if err != nil {
		return fmt.Errorf("incorrect product name: %v", err)
	}

	offer.ProductId = productId
	return nil
}

func (o *ProductOfferUseCase) GetOfferForProduct(ctx context.Context, productName string) ([]prdtoffer.Offer, error) {
	productId, err := o.prod.GetIDByName(ctx, productName)
	if err != nil {
		return nil, fmt.Errorf("incorrect product name: %v", err)
	}

	modelOffer, err := o.rep.GetOfferForProduct(ctx, productId)
	if err != nil {
		return nil, err
	}

	return o.arrayModToOffer(modelOffer)
}

func (o *ProductOfferUseCase) arrayModToOffer(mod []models.Offer) ([]prdtoffer.Offer, error) {
	result := make([]prdtoffer.Offer, 0)
	for _, val := range mod {
		tmp, err := o.modToOffer(val)
		if err != nil {
			return nil, err
		}

		result = append(result, *tmp)
	}

	return result, nil
}

func (o *ProductOfferUseCase) modToOffer(mod models.Offer) (*prdtoffer.Offer, error) {
	result, err := &prdtoffer.Offer{}, error(nil)

	result.ProductName, err = o.prod.GetNameById(mod.ProductId)
	if err != nil {
		return nil, err
	}

	result.ProviderName, err = o.prov.GetNameById(mod.ProviderId)
	if err != nil {
		return nil, err
	}

	result.Id = mod.Id
	result.Cost = mod.Cost

	return result, nil
}

func (o *ProductOfferUseCase) GetOffersOfProvider(ctx context.Context, providerName string) ([]prdtoffer.Offer, error) {
	providerId, err := o.prov.GetIDByName(ctx, providerName)
	if err != nil {
		return nil, fmt.Errorf("incorrect provider name: %v", err)
	}

	offers, err := o.rep.GetOfferOfProvider(ctx, providerId)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return o.arrayModToOffer(offers)
}

func (o *ProductOfferUseCase) UpdateCost(ctx context.Context, providerName, productName string, cost float32) error {
	providerId, err := o.prov.GetIDByName(ctx, providerName)
	if err != nil {
		return fmt.Errorf("incorrect provider name: %v", err)
	}

	productId, err := o.prod.GetIDByName(ctx, productName)
	if err != nil {
		return fmt.Errorf("incorrect product name: %v", err)
	}

	if err := o.rep.UpdateCost(ctx, providerId, productId, cost); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func (o *ProductOfferUseCase) Delete(ctx context.Context, productName, providerName string) error {
	providerId, err := o.prov.GetIDByName(ctx, providerName)
	if err != nil {
		return fmt.Errorf("incorrect provider name: %v", err)
	}

	productId, err := o.prod.GetIDByName(ctx, productName)
	if err != nil {
		return fmt.Errorf("incorrect product name: %v", err)
	}

	if err := o.rep.Delete(ctx, providerId, productId); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func (o *ProductOfferUseCase) GetById(id int) (*prdtoffer.Offer, error) {
	offer_base, err := o.rep.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return o.modToOffer(*offer_base)
}
