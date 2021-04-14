package app

import (
	"context"
	"cw/models"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JWTAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isNoAuthPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		tokenString, status := get_token(r)
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		tokenInfo, statusCode := decodingJWT(tokenString)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tokenInfo.Login)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func isNoAuthPath(path string) bool {
	for _, noAuthPath := range []string{"/login", "/register", "/refresh_token"} {
		if path == noAuthPath {
			return true
		}
	}

	return false
}

func decodingJWT(tokenString string) (*models.TokenInfo, int) {
	tk := &models.TokenInfo{}
	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEYPASSWORD")), nil
	})

	if err != nil {
		log.Println("decodintJWT: ", err)
		return nil, http.StatusForbidden
	}

	return tk, checkToken(token)
}

func checkToken(token *jwt.Token) int {
	if !token.Valid {
		return http.StatusForbidden
	}

	return http.StatusOK
}

func get_token(r *http.Request) (string, int) {
	tokenHeader := r.Header.Get("Authorization")
	if isEmptyToken(tokenHeader) {
		return "", http.StatusUnauthorized
	}

	return getTokenFromHeaderValue(tokenHeader)
}

func isEmptyToken(token string) bool {
	return token == ""
}

func getTokenFromHeaderValue(tokenHeaderValue string) (string, int) {
	splitted := strings.Split(tokenHeaderValue, " ")
	if isInvalidFormatOfToken(splitted) {
		return "", http.StatusBadRequest
	}
	return splitted[1], http.StatusOK
}

func isInvalidFormatOfToken(headerWord []string) bool {
	return len(headerWord) != 2
}

var LogNewConnection = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isNoAuthPath(r.URL.Path) {
			user := r.Context().Value("user")
			log.Printf("Connection user: %v", user)
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})

}
