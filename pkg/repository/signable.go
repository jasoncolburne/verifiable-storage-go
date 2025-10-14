package repository

import (
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SignableRepository[T primitives.SignableAndRecordable] struct {
	VerifiableRepository[T]

	signingKey interfaces.SigningKey
}

func (s SignableRepository[T]) CreateVersion(record T) error {
	prepareSignableRecord(record, s.noncer, s.signingKey)

	if s.write {
		// write to data store
	}

	return nil
}
