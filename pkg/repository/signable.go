package repository

import (
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SignableRepository[T primitives.SignableAndRecordable] struct {
	VerifiableRepository[T]

	signingKey interfaces.SigningKey
}

func NewSignableRepository[T primitives.SignableAndRecordable](store data.Store, write bool, noncer interfaces.Noncer, signingKey interfaces.SigningKey) *SignableRepository[T] {
	return &SignableRepository[T]{
		VerifiableRepository: VerifiableRepository[T]{
			store:  store,
			noncer: noncer,

			write: write,
		},

		signingKey: signingKey,
	}
}

func (s SignableRepository[T]) CreateVersion(record T) error {
	prepareSignableRecord(record, s.noncer, s.signingKey)

	if s.write {
		// write to data store
	}

	return nil
}
