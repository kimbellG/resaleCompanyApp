package usecase

import (
	"context"
	"cw/auth"
	"cw/models"
	"fmt"
)

type AdminUseCase struct {
	repo auth.AdminRepository
}

func NewAdminUseCase(repo auth.AdminRepository) *AdminUseCase {
	return &AdminUseCase{
		repo: repo,
	}
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

func (a *AdminUseCase) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := a.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return users, nil
}

func (a *AdminUseCase) GetUser(ctx context.Context, username string) (*models.User, error) {
	user, err := a.repo.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return user, nil
}

func (a *AdminUseCase) UpdateUser(ctx context.Context, username string, key, value string) error {
	if key == "status" {
		return fmt.Errorf("incorrect update's field")
	}

	if err := a.repo.UpdateUser(ctx, username, key, value); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func (a *AdminUseCase) DeleteUser(ctx context.Context, username string) error {
	if err := a.repo.DeleteUser(ctx, username); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}
