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

// type KVStore interface {
// 	Get(ctx context.Context, key string) (any, bool, error)
// 	Set(ctx context.Context, key string, v any) error
// 	Delete(ctx context.Context, key string) error
// }

type Store interface {
	Sql() SQLStore // May be an sql.Tx or sql.DB
	// Kv() KVStore

	Placeholders(count int) []string
}
