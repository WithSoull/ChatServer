package config

import "github.com/joho/godotenv"

func Load(path string) error {
  if err := godotenv.Load(path); err != nil {
    return err
  }
  return nil
}

type GRPCCongif interface {
  Address() string
}

type PGCongif interface {
  DSN() string
}
