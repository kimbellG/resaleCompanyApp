package usecase

import (
	"context"
	"cw/auth"
	"cw/logger"
	"cw/models"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type AuthUseCase struct {
	userRepo auth.UserRepository
	tokenKey []byte
}

func NewAuthUseCase(rep auth.UserRepository, passwordKey []byte) *AuthUseCase {
	return &AuthUseCase{
		userRepo: rep,
		tokenKey: passwordKey,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username, password, name string) error {
	new_user := &models.User{
		Login:    username,
		Password: password,
		Name:     name,
		Access:   "",
		Status:   false,
	}

	return a.userRepo.CreateUser(ctx, new_user)
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (string, error) {
	user, err := a.userRepo.GetUser(ctx, username, password)
	if err != nil {
		return "", err
	}

	if !user.Status {
		return "", errors.New("Profile isn't activate")
	}

	return a.GenerateToken(user), nil
}

func NewTokenInfo(user *models.User) *models.TokenInfo {
	return &models.TokenInfo{
		AccessProfile:  user.Access,
		StandardClaims: jwt.StandardClaims{},
	}
}

func (a *AuthUseCase) GenerateToken(user *models.User) string {
	tk := NewTokenInfo(user)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString(a.tokenKey)
	if err != nil {
		logger.AssertMessage(
			map[string]interface{}{"action": "generate token"},
			fmt.Sprintf("invalid procces of creating token: %v", err))
	}

	return tokenString
}
