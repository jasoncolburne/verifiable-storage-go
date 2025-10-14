package repository

import (
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

func (v VerifiableRepository[T]) CreateVersion(record T) error {
	prepareVerifiableRecord(record, v.noncer)

	if v.write {
		// write to data store
	}

	return nil
}
