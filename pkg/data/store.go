package data

import (
	"context"
	"database/sql"
)

type SQLStore interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// type KVStore interface {
// 	Get(ctx context.Context, key string) (any, bool, error)
// 	Set(ctx context.Context, key string, v any) error
// 	Delete(ctx context.Context, key string) error
// }

type Store interface {
	Sql() SQLStore // May be an sql.Tx or sql.DB
	// Kv() KVStore

	Placeholder() string
}
