package auth

import (
	"context"
	"cw/models"
)

type UseCase interface {
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, assessToken string) (*models.User, error)
}
