package models

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenAnswer struct {
	CurrentToken string `json:"current_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInfo struct {
	AccessProfile string
	Sec           time.Time
	TypeToken     byte
	jwt.StandardClaims
}

func (user *AuthorizationUserInformation) GenerateToken() *TokenAnswer {
	tkClaimsCurrent := &TokenInfo{AccessProfile: user.AccessProfile, Sec: time.Now().Add(time.Hour), TypeToken: 'c'}
	tkClaimsRefresh := &TokenInfo{AccessProfile: user.AccessProfile, Sec: time.Now().Add(time.Hour * 24 * 30), TypeToken: 'r'}

	return &TokenAnswer{tkClaimsCurrent.GenerateToken(), tkClaimsRefresh.GenerateToken()}
}

func (tk *TokenInfo) GenerateToken() string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("KEYPASSWORD")))
	if err != nil {
		log.Fatalln(err)
	}

	return tokenString
}

func (tokens *TokenAnswer) EncodingTokensInHttpBody(w http.ResponseWriter) (http.ResponseWriter, error) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		return w, err
	}

	return w, nil
}
