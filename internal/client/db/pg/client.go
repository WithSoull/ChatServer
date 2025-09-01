package pg

import (
	"context"

	"github.com/pkg/errors"

	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgClient struct {
	masterDBC db.DB
}

func NewPGClient(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err.Error())
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
