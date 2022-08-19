package core

import "github.com/golang-jwt/jwt"

type AppClaim struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}
