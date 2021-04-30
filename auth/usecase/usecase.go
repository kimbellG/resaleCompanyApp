package usecase

import (
	"context"
	"cw/auth"
	"cw/logger"
	"cw/models"
	"errors"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

type AuthUseCase struct {
	userRepo *auth.UserRepository
	tokenKey []byte
}

func NewAuthUseCase(rep *auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: rep,
		tokenKey: []byte(os.Getenv("KEYPASSWORD")),
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username, password, name, access string) error {
	new_user := models.User{
		Login:    username,
		Password: password,
		Name:     name,
		Access:   access,
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

type TokenInfo struct {
	AccessProfile string
	jwt.StandardClaims
}

func NewTokenInfo(user *models.User) *TokenInfo {
	return &TokenInfo{
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
