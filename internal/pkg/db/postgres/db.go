package postgres

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Connection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type DB struct {
	conn Connection
}

func NewDB(conn Connection) *DB {
	return &DB{conn: conn}
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	conn := db.getConn(ctx)
	return conn.Exec(ctx, sql, args...)
}

func (db *DB) Get(ctx context.Context, dst any, sql string, args ...any) error {
	conn := db.getConn(ctx)
	return pgxscan.Get(ctx, conn, dst, sql, args...)
}

func (db *DB) Select(ctx context.Context, dst any, sql string, args ...any) error {
	conn := db.getConn(ctx)
	return pgxscan.Select(ctx, conn, dst, sql, args...)
}

func (db *DB) getConn(ctx context.Context) Connection {
	conn, ok := ctx.Value(transactionKey).(Connection)
	if !ok {
		return db.conn
	}
	return conn
}
