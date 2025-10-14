package repository

import "github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"

type SearchableRepository[T primitives.SearchableAndRecordable] struct {
	VerifiableRepository[T]
}

func (s SearchableRepository[T]) CreateVersion(record T) error {
	prepareSearchableRecord(record, s.noncer)

	if s.write {
		// write to data store
	}

	return nil
}
