package auth

import (
	"gopkg.in/yaml.v3"
)

type AuthRoot struct {
	Auth AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	JwtAuth *JwtAuthConfig `yaml:"jwt-auth"`
}

type JwtAuthConfig struct {
	Secret     string `yaml:"secret"`
	AuthHeader string `yaml:"auth-header"`
	Claims     Claims `json:"claims"`
}

type Claims struct {
	UserId string `yaml:"user-id"`
	Role   string `yaml:"role"`
}

func (j JwtAuthConfig) toProxy() JwtAuthProxy {
	return JwtAuthProxy{
		secret:     j.Secret,
		authHeader: j.AuthHeader,
		Claims: Claims{
			UserId: j.Claims.UserId,
			Role:   j.Claims.Role,
		},
	}
}

func SetUpAuth(data []byte) error {
	var config AuthRoot
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	auth := config.Auth.JwtAuth
	if auth != nil {
		proxy := auth.toProxy()
		Save(&proxy)
	}
	return nil
}
