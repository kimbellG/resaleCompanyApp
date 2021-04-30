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
	AccessProfile string
	jwt.StandardClaims
}
