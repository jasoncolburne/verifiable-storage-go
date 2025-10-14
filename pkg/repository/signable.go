package repository

import (
	"context"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SignableRepository[T primitives.SignableAndRecordable] struct {
	VerifiableRepository[T]

	signingKey           interfaces.SigningKey
	verificationKeyStore interfaces.VerificationKeyStore
}

func NewSignableRepository[T primitives.SignableAndRecordable](
	store data.Store,
	write bool,
	noncer interfaces.Noncer,
	signingKey interfaces.SigningKey,
	verificationKeyStore interfaces.VerificationKeyStore,
) *SignableRepository[T] {
	return &SignableRepository[T]{
		VerifiableRepository: VerifiableRepository[T]{
			store:  store,
			noncer: noncer,

			write: write,
		},

		signingKey:           signingKey,
		verificationKeyStore: verificationKeyStore,
	}
}

func (r SignableRepository[T]) CreateVersion(ctx context.Context, record T) error {
	if err := prepareSignedRecord(record, r.noncer, r.signingKey); err != nil {
		return err
	}

	if r.write {
		if err := r.insertRecord(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func (r SignableRepository[T]) GetById(ctx context.Context, record T, id string) error {
	if err := r.getRecord(ctx, record, id); err != nil {
		return err
	}

	if err := verifySignedRecord(record, r.verificationKeyStore); err != nil {
		return err
	}

	return nil
}
