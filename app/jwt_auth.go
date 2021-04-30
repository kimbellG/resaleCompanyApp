package app

import (
	"context"
	"cw/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type access string

const accessProfile = access("access")

var JWTAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
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

		ctx := context.WithValue(r.Context(), accessProfile, tokenInfo.AccessProfile)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func isPublicPath(path string) bool {
	for _, noAuthPath := range []string{"/sign-up", "/register", "/refresh_token"} {
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
		if !isPublicPath(r.URL.Path) {
			user := r.Context().Value(accessProfile)
			log.Printf("Connection user: %v", user)
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})

}

var CheckAccessRight = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			log.Println("is public path: don't check access right")
			next.ServeHTTP(w, r)
			return
		}

		accessRight, err := readAccessRightFromFile()
		if err != nil {
			log.Fatalf("error: read access right from file: %v", err)
		}

		AccessProfile := r.Context().Value(accessProfile)
		AccessPaths, ok := accessRight[fmt.Sprintf("%v", AccessProfile)]
		if !ok {
			log.Println("Incorrect access profile:", AccessProfile)
			http.Error(w, "Incorrect access profile", http.StatusForbidden)
			return
		}

		if err := checkURLPATH(AccessPaths, r.URL.Path); err != nil {
			log.Println(r.URL.Path, ":", err)
			http.Error(w, fmt.Sprint(err), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func readAccessRightFromFile() (map[string][]string, error) {
	accessRight := map[string][]string{}
	file, err := os.Open("access_right.json")
	if err != nil {
		return nil, errors.Errorf("file access doesn't open")
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&accessRight); err != nil {
		return nil, errors.Errorf("invalid json format")
	}

	return accessRight, nil
}

func checkURLPATH(AccessPaths []string, URLPath string) error {
	for _, path := range AccessPaths {
		if URLPath == path {
			return nil
		}
	}

	return errors.New("resource is not available")
}
