package repository

import (
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SearchableRepository[T primitives.SearchableAndRecordable] struct {
	VerifiableRepository[T]
}

func NewSearchableRepository[T primitives.SearchableAndRecordable](store data.Store, write bool, noncer interfaces.Noncer) *SearchableRepository[T] {
	return &SearchableRepository[T]{
		VerifiableRepository: VerifiableRepository[T]{
			store:  store,
			noncer: noncer,

			write: write,
		},
	}
}

func (s SearchableRepository[T]) CreateVersion(record T) error {
	prepareSearchableRecord(record, s.noncer)

	if s.write {
		// write to data store
	}

	return nil
}
