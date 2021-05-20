package models

import jwt "github.com/dgrijalva/jwt-go"

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	Access   string `json:"access"`
	Name     string `json:"name"`
}

type TokenInfo struct {
	Login         string
	AccessProfile string
	jwt.StandardClaims
}
