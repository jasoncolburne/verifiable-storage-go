package repository

import (
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SignableAndSearchableRepository[T primitives.SignableAndSearchableAndRecordable] struct {
	SignableRepository[T]
}

func (s SignableAndSearchableRepository[T]) CreateVersion(record T) error {
	prepareSignableSearchableRecord(record, s.noncer, s.signingKey)

	if s.write {
		// write to data store
	}

	return nil
}
