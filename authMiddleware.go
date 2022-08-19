package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

func shouldCheckToken(pagesPass []string, route string) bool {
	for _, p := range pagesPass {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func AuthJWT(pagesPass []string, s Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(pagesPass, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenStr, &AppClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetTokenJWT(userId int, s Server) (string, error) {
	claims := AppClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
	if err != nil {
		return "", errors.New("Unauthorized")
	}
	return tokenString, nil
}

func CompareHashAndPasswordJWT(pass1, pass2 string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(pass1), []byte(pass2)); err != nil {
		return errors.New("Data invalid")
	}
	return nil
}

func GenerateHashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("Password required")
	}
	if len(password) <= 3 {
		return "", errors.New("Password length invalid")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), HASH_COST)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}
