package repository

import (
	"context"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type VerifiableRepository[T primitives.VerifiableAndRecordable] struct {
	store  data.Store
	noncer interfaces.Noncer

	// enable writes (disabled for admin dry-run commands for instance)
	write bool
}

func NewVerifiableRepository[T primitives.VerifiableAndRecordable](store data.Store, write bool, noncer interfaces.Noncer) *VerifiableRepository[T] {
	return &VerifiableRepository[T]{
		store:  store,
		noncer: noncer,

		write: write,
	}
}

func (r VerifiableRepository[T]) CreateVersion(ctx context.Context, record T) error {
	if err := prepareVerifiableRecord(record, r.noncer); err != nil {
		return err
	}

	if r.write {
		r.store.Sql().ExecContext(
			ctx,
			``,
		)
	}

	return nil
}
