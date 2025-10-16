package examples

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewInMemorySQLiteStore() (*SQLiteStore, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db.Unsafe(), // the unsafe here allows us to gracefully ignore computed columns
		tx: nil,
	}, nil
}

func (s SQLiteStore) Sql() data.SQLStore {
	if s.tx == nil {
		return s.db
	} else {
		return s.tx
	}
}

func (s *SQLiteStore) BeginTransaction(ctx context.Context, opts *sql.TxOptions) error {
	if s.tx != nil {
		return fmt.Errorf("transaction in progress")
	}

	var err error
	s.tx, err = s.db.BeginTxx(ctx, opts)
	if err != nil {
		s.tx = nil
		return err
	}

	return nil
}

func (s *SQLiteStore) CommitTransaction() error {
	if s.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	if err := s.tx.Commit(); err != nil {
		return err
	}

	s.tx = nil

	return nil
}

func (s *SQLiteStore) RollbackTransaction() error {
	if s.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	if err := s.tx.Rollback(); err != nil {
		return err
	}

	s.tx = nil

	return nil
}
