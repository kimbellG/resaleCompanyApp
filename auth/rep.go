package auth

import (
	"context"
	"cw/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username, password string) (*models.User, error)
}

type AdminRepository interface {
	OnOffUser(ctx context.Context, username string, status bool) error
	SetAccessProfile(ctx context.Context, username string, access string) error

	GetUser(ctx context.Context, username string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)

	UpdateUser(ctx context.Context, username string, key, value string) error
	DeleteUser(ctx context.Context, username string) error
}
