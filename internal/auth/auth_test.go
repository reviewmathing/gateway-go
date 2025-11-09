package auth

import (
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJwtAuthProxy(t *testing.T) {
	secretValue := "testsecrettestsecrettestsecrettestsecrettestsecrettestsecrettestsecret"
	proxy := JwtAuthProxy{
		secret:     secretValue,
		authHeader: "authentication",
		Claims: Claims{
			UserId: "userId",
			Role:   "role",
		},
	}

	claims := jwt.MapClaims{}
	claims["userId"] = "testUser"
	claims["role"] = []string{"ADMIN", "USER"}
	withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := withClaims.SignedString([]byte(secretValue))
	if err != nil {
		t.Fatal("token create fail ", err)
	}

	request := http.Request{
		Header: map[string][]string{},
	}
	request.Header.Set("authentication", signedString)

	err = proxy.Handle(&request)
	if err != nil {
		t.Fatal("Auth handle fail ", err)
	}

	if request.Header.Get("X-USER-ID") == "" {
		t.Fatal("X-USER-ID is empty")
	}

	if request.Header.Get("X-USER-ROLE") == "" {
		t.Fatal("X-USER-ID is empty")
	}
}
