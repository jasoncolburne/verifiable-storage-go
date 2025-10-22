package data

import (
	"context"
	"database/sql"
)

type SQLStore interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type Store interface {
	Sql() SQLStore // May be an sql.Tx or sql.DB

	BeginTransaction(ctx context.Context, opts *sql.TxOptions) error
	CommitTransaction() error
	RollbackTransaction() error

	ReplacePlaceholders(query string) string
}
