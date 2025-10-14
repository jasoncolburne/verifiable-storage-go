package repository

import (
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SignableAndSearchableRepository[T primitives.SignableAndSearchableAndRecordable] struct {
	SignableRepository[T]
}

func NewSignableAndSearchableRepository[T primitives.SignableAndSearchableAndRecordable](store data.Store, write bool, noncer interfaces.Noncer, signingKey interfaces.SigningKey) *SignableAndSearchableRepository[T] {
	return &SignableAndSearchableRepository[T]{
		SignableRepository: SignableRepository[T]{
			VerifiableRepository: VerifiableRepository[T]{
				store:  store,
				noncer: noncer,

				write: write,
			},

			signingKey: signingKey,
		},
	}
}

func (s SignableAndSearchableRepository[T]) CreateVersion(record T) error {
	prepareSignableSearchableRecord(record, s.noncer, s.signingKey)

	if s.write {
		// write to data store
	}

	return nil
}
