package config

import "github.com/joho/godotenv"

func Load(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}
	return nil
}

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}
