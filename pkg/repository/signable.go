package repository

import (
	"encoding/json"
	"fmt"

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

func (r SignableRepository[T]) CreateVersion(record T) error {
	prepareSignedRecord(record, r.noncer, r.signingKey)

	if r.write {
		// write to data store
		bytes, err := json.Marshal(record)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", bytes)
	}

	return nil
}
