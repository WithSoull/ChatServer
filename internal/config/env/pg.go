package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/WithSoull/ChatServer/internal/config"
)

const (
	dsnEnvName      = "PG_DSN"
	hostEnvName     = "PG_HOST"
	portEnvName     = "PG_PORT_INNER"
	dbEnvName       = "PG_DATABASE_NAME"
	userEnvName     = "PG_USER"
	passwordEnvName = "PG_PASSWORD"
	sslModeEnvName  = "PG_SSL_MODE"
)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (config.PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		host := os.Getenv(hostEnvName)
		port := os.Getenv(portEnvName)
		dbname := os.Getenv(dbEnvName)
		user := os.Getenv(userEnvName)
		password := os.Getenv(passwordEnvName)
		sslmode := os.Getenv(sslModeEnvName)

		if len(host) == 0 ||
			len(port) == 0 ||
			len(dbname) == 0 ||
			len(user) == 0 ||
			len(password) == 0 {
			return nil, errors.New("missing required environment variables for database connection (Tip: you can just define *PG_DSN in environment)")
		}

		if len(sslmode) == 0 {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			host, port, dbname, user, password, sslmode,
		)
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
