package usecase

import (
	"context"
	"cw/models"
	"cw/provider"
)

type ProviderUseCase struct {
	rep provider.Repository
}

func NewProviderUseCase(rep_pr provider.Repository) *ProviderUseCase {
	return &ProviderUseCase{
		rep: rep_pr,
	}
}

func (p *ProviderUseCase) AddProvider(ctx context.Context, pr *provider.Provider) error {
	modProvider := prToModels(pr)

	return p.rep.AddProvider(ctx, modProvider)
}

func prToModels(pr *provider.Provider) *models.Provider {
	return &models.Provider{
		VendorCode:     pr.VendorCode,
		Name:           pr.Name,
		UNP:            pr.UNP,
		TermsOfPayment: pr.TermsOfPayment,
		PhoneNumber:    pr.PhoneNumber,
		Address:        pr.Address,
		Email:          pr.Email,
		WebSite:        pr.WebSite,
	}
}

func (p *ProviderUseCase) GetProviders(ctx context.Context) ([]provider.Provider, error) {
	modProviders, err := p.rep.GetProviders(ctx)
	if err != nil {
		return nil, err
	}

	result := new([]provider.Provider)
	for _, mod := range modProviders {
		*result = append(*result, *modToPr(&mod))
	}

	return *result, nil
}

func modToPr(mod *models.Provider) *provider.Provider {
	return &provider.Provider{
		VendorCode:     mod.VendorCode,
		Name:           mod.Name,
		UNP:            mod.UNP,
		TermsOfPayment: mod.TermsOfPayment,
		PhoneNumber:    mod.PhoneNumber,
		Address:        mod.Address,
		Email:          mod.Email,
		WebSite:        mod.WebSite,
	}
}

func (p *ProviderUseCase) UpdateProvider(ctx context.Context, code int, fields map[string]interface{}) error {
	return p.rep.UpdateProvider(ctx, code, fields)
}

func (p *ProviderUseCase) DeleteProvider(ctx context.Context, code int) error {
	return p.rep.DeleteProvider(ctx, code)
}
