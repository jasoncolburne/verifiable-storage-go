package repository

import (
	"context"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type SignableRepository[T primitives.SignableAndRecordable] struct {
	VerifiableRepository[T]

	signingKey           interfaces.SigningKey
	verificationKeyStore interfaces.VerificationKeyStore
}

// pass a nil noncer to omit nonces
func NewSignableRepository[T primitives.SignableAndRecordable](
	store data.Store,
	write bool,
	timestamp bool,
	noncer interfaces.Noncer,
	signingKey interfaces.SigningKey,
	verificationKeyStore interfaces.VerificationKeyStore,
) *SignableRepository[T] {
	return &SignableRepository[T]{
		VerifiableRepository: VerifiableRepository[T]{
			store:  store,
			noncer: noncer,

			write:     write,
			timestamp: timestamp,
		},

		signingKey:           signingKey,
		verificationKeyStore: verificationKeyStore,
	}
}

func (r SignableRepository[T]) CreateVersion(ctx context.Context, record T) error {
	if err := r.prepareSignedRecord(record); err != nil {
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
	if err := r.getRecordById(ctx, record, id); err != nil {
		return err
	}

	if err := r.verifySignedRecord(record); err != nil {
		return err
	}

	return nil
}

func (r SignableRepository[T]) GetLatestByPrefix(ctx context.Context, record T, prefix string) error {
	if err := r.getLatestRecordByPrefix(ctx, record, prefix); err != nil {
		return err
	}

	if err := r.verifySignedRecord(record); err != nil {
		return err
	}

	return nil
}

func (r SignableRepository[T]) ListByPrefix(ctx context.Context, records *[]T, prefix string) error {
	if err := r.listRecordsByPrefix(ctx, records, prefix); err != nil {
		return err
	}

	for _, record := range *records {
		if err := r.verifySignedRecord(record); err != nil {
			return err
		}
	}

	return nil
}

func (r SignableRepository[T]) Get(
	ctx context.Context,
	record T,
	condition data.ClauseOrExpression,
	order data.Ordering,
) error {
	if err := r.get(ctx, record, condition, order); err != nil {
		return err
	}

	if err := r.verifySignedRecord(record); err != nil {
		return err
	}

	return nil
}

func (r SignableRepository[T]) Select(
	ctx context.Context,
	records *[]T,
	condition data.ClauseOrExpression,
	order data.Ordering,
	limit *uint,
) error {
	if err := r._select(ctx, records, condition, order, limit); err != nil {
		return err
	}

	for _, record := range *records {
		if err := r.verifySignedRecord(record); err != nil {
			return err
		}
	}

	return nil
}

func (r SignableRepository[T]) ListLatestByPrefix(
	ctx context.Context,
	records *[]T,
	condition data.ClauseOrExpression,
	order data.Ordering,
	limit *uint,
) error {
	if err := r.selectLatestByPrefix(ctx, records, condition, order, limit); err != nil {
		return err
	}

	for _, record := range *records {
		if err := r.verifySignedRecord(record); err != nil {
			return err
		}
	}

	return nil
}

// helpers

func (r SignableRepository[T]) prepareSignedRecord(record T) error {
	if err := algorithms.Sign(record, r.signingKey, func() error {
		return r.prepareVerifiableRecord(record)
	}); err != nil {
		return err
	}

	return nil
}

func (r SignableRepository[T]) verifySignedRecord(record T) error {
	if err := algorithms.VerifySignature(record, r.verificationKeyStore); err != nil {
		return err
	}

	if err := r.verifyRecord(record); err != nil {
		return err
	}

	return nil
}
