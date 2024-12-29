package postgres

import (
	"context"

	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/jackc/pgx/v5"
)

type TransactionableConnection interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type TransactionManager struct {
	conn TransactionableConnection
}

func (tm *TransactionManager) Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel: pgxTxIsolation(isolation),
	}

	return pgx.BeginTxFunc(ctx, tm.conn, opts, func(tx pgx.Tx) error {
		txCtx := context.WithValue(ctx, transactionKey, tx)

		err := f(txCtx)
		if err != nil {
			return err
		}

		return nil
	})
}

func pgxTxIsolation(i db.TxIsolation) pgx.TxIsoLevel {
	switch i {
	case db.TxIsolationReadUncommitted:
		return pgx.ReadCommitted
	case db.TxIsolationReadCommitted:
		return pgx.ReadCommitted
	case db.TxIsolationRepeatableRead:
		return pgx.RepeatableRead
	case db.TxIsolationSerializable:
		return pgx.Serializable
	default:
		return ""
	}
}
