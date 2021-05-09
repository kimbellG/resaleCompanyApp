package models

import jwt "github.com/dgrijalva/jwt-go"

type User struct {
	Login    string
	Password string
	Status   bool
	Access   string
	Name     string
}

type TokenInfo struct {
	Login         string
	AccessProfile string
	jwt.StandardClaims
}
