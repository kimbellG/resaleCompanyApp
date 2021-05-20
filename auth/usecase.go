package auth

import (
	"context"
	"cw/models"
)

type UseCase interface {
	SignUp(ctx context.Context, username, password, name string) error
	SignIn(ctx context.Context, username, password string) (string, error)
}

type AdminUseCase interface {
	ConfirmUser(ctx context.Context, username string, accessProfile string) error
	DisableUser(ctx context.Context, username string) error

	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)

	UpdateUser(ctx context.Context, username string, key, value string) error
	DeleteUser(ctx context.Context, username string) error
}
