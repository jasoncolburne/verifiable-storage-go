package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type VerifiableRepository[T primitives.VerifiableAndRecordable] struct {
	store  data.Store
	noncer interfaces.Noncer

	// enable writes (disabled for admin dry-run commands for instance)
	write bool
}

func NewVerifiableRepository[T primitives.VerifiableAndRecordable](store data.Store, write bool, noncer interfaces.Noncer) *VerifiableRepository[T] {
	return &VerifiableRepository[T]{
		store:  store,
		noncer: noncer,

		write: write,
	}
}

func (r VerifiableRepository[T]) CreateVersion(ctx context.Context, record T) error {
	if err := prepareVerifiableRecord(record, r.noncer); err != nil {
		return err
	}

	if r.write {
		if err := r.insertRecord(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func (r VerifiableRepository[T]) GetById(ctx context.Context, record T, id string) error {
	if err := r.getRecord(ctx, record, id); err != nil {
		return err
	}

	if err := verifyRecord(record); err != nil {
		return err
	}

	return nil
}

func (r VerifiableRepository[T]) GetLatestByPrefix(ctx context.Context, record T, prefix string) error {
	if err := r.getLatestRecordByPrefix(ctx, record, prefix); err != nil {
		return err
	}

	if err := verifyRecord(record); err != nil {
		return err
	}

	return nil
}

// helpers

func (r VerifiableRepository[T]) insertRecord(ctx context.Context, record T) error {
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

	return nil
}

func (r VerifiableRepository[T]) getRecord(ctx context.Context, record T, id string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %s", record.TableName(), r.store.Placeholder())

	if err := r.store.Sql().GetContext(ctx, record, query, id); err != nil {
		return err
	}

	return nil
}

func (r VerifiableRepository[T]) getLatestRecordByPrefix(ctx context.Context, record T, prefix string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE prefix = %s ORDER BY sequence_number DESC LIMIT 1", record.TableName(), r.store.Placeholder())

	if err := r.store.Sql().GetContext(ctx, record, query, prefix); err != nil {
		return err
	}

	return nil
}
