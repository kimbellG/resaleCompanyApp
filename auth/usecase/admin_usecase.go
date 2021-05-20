package usecase

import (
	"context"
	"cw/auth"
	"fmt"
)

type AdminUseCase struct {
	repo auth.AdminRepository
}

func (a *AdminUseCase) ConfirmUser(ctx context.Context, username string, accessProfile string) error {
	if err := a.repo.OnOffUser(ctx, username, true); err != nil {
		return fmt.Errorf("activate user: %v", err)
	}

	if err := a.repo.SetAccessProfile(ctx, username, accessProfile); err != nil {
		return fmt.Errorf("set access profile: %v", err)
	}

	return nil
}

func (a *AdminUseCase) DisableUser(ctx context.Context, username string) error {
	if err := a.repo.OnOffUser(ctx, username, false); err != nil {
		return fmt.Errorf("disable user: %v", err)
	}

	return nil
}
