package txctx

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/contextx"
	"github.com/jackc/pgx/v4"
)

const txKey contextx.CtxKey = "tx"

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func ExtractTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return nil, false
	}
	return tx, true
}
