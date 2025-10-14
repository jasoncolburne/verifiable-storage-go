package repository

import (
	"context"
	"fmt"
	"strings"

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

func (r SignableRepository[T]) CreateVersion(ctx context.Context, record T) error {
	if err := prepareSignedRecord(record, r.noncer, r.signingKey); err != nil {
		return err
	}

	if r.write {
		fieldNames := getFieldNames(record)
		innerFields := strings.Join(fieldNames, ", ")
		innerValues := strings.Join(fieldNames, ", :")

		fmt.Printf("%s\n", innerFields)
		fmt.Printf("%s\n", innerValues)

		// write to data store
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", record.TableName(), innerFields, innerValues)

		_, err := r.store.Sql().NamedExecContext(
			ctx,
			query,
			record,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r SignableRepository[T]) GetById(ctx context.Context, t T, id string) error {
	query := fmt.Sprintf("SELECT * from %s where id = %s", t.TableName(), r.store.Placeholder())

	if err := r.store.Sql().GetContext(ctx, t, query, id); err != nil {
		return err
	}

	return nil
}
