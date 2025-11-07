package env

import (
	"github.com/caarlos0/env/v11"
)

type jwtEnvConfig struct {
	RefreshSecret string `env:"REFRESH_TOKEN_SECRET,notEmpty"`
	AccessSecret  string `env:"ACCESS_TOKEN_SECRET,notEmpty"`
}

type jwtConfig struct {
	raw jwtEnvConfig
}

func NewJWTConfig() (*jwtConfig, error) {
	var raw jwtEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &jwtConfig{raw: raw}, nil
}

func (c *jwtConfig) RefreshTokenSecretKey() string {
	return c.raw.RefreshSecret
}

func (c *jwtConfig) AccessTokenSecretKey() string {
	return c.raw.AccessSecret
}
