package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthType string

const (
	NONE AuthType = "NONE"
	JWT  AuthType = "JWT"
)

func ParseAuthType(value string) AuthType {
	switch strings.ToUpper(value) {
	case "JWT":
		return JWT
	default:
		return NONE
	}
}

type AuthError error

func NewAuthError(message string) AuthError {
	return errors.New(message)
}

type ProxyType AuthType

type AuthProxy interface {
	Handle(*http.Request) error
	GetType() ProxyType
}

type JwtAuthProxy struct {
	secret     string
	authHeader string
	Claims     Claims
}

func (j *JwtAuthProxy) Handle(r *http.Request) error {
	authValue := r.Header.Get(j.authHeader)
	if authValue == "" {
		return NewAuthError("Authentication header is empty")
	}
	claims := jwt.MapClaims{}
	key := []byte(j.secret)
	token, err := jwt.ParseWithClaims(authValue, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil || !token.Valid {
		return err
	}

	userId := claims[j.Claims.UserId].(string)
	rolesInterface := claims["role"].([]interface{}) // interface{} 슬라이스로 가져오기

	roles := make([]string, len(rolesInterface))
	if userId == "" || roles == nil || len(roles) == 0 {
		return NewAuthError("")
	}

	updateRequest(userId, roles, r)
	return nil
}

func (j *JwtAuthProxy) GetType() ProxyType {
	return ProxyType(JWT)
}

func updateRequest(userId string, roles []string, r *http.Request) {
	r.Header.Set("X-User-Id", userId)
	r.Header.Set("X-User-Role", strings.Join(roles, ","))
}
