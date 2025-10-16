package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type VerifiableRepository[T primitives.VerifiableAndRecordable] struct {
	store  data.Store
	noncer interfaces.Noncer

	// enable writes (disabled for admin dry-run commands for instance)
	write     bool
	timestamp bool
}

// pass a nil noncer to omit nonces
func NewVerifiableRepository[T primitives.VerifiableAndRecordable](
	store data.Store,
	write bool,
	timestamp bool,
	noncer interfaces.Noncer,
) *VerifiableRepository[T] {
	return &VerifiableRepository[T]{
		store:  store,
		noncer: noncer,

		write:     write,
		timestamp: timestamp,
	}
}

func (r VerifiableRepository[T]) CreateVersion(ctx context.Context, record T) error {
	if err := r.prepareVerifiableRecord(record); err != nil {
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
	if err := r.getRecordById(ctx, record, id); err != nil {
		return err
	}

	if err := r.verifyRecord(record); err != nil {
		return err
	}

	return nil
}

func (r VerifiableRepository[T]) GetLatestByPrefix(ctx context.Context, record T, prefix string) error {
	if err := r.getLatestRecordByPrefix(ctx, record, prefix); err != nil {
		return err
	}

	if err := r.verifyRecord(record); err != nil {
		return err
	}

	return nil
}

func (r VerifiableRepository[T]) ListByPrefix(ctx context.Context, records *[]T, prefix string) error {
	if err := r.listRecordsByPrefix(ctx, records, prefix); err != nil {
		return err
	}

	for _, record := range *records {
		if err := r.verifyRecord(record); err != nil {
			return err
		}
	}

	return nil
}

// helpers

func (r VerifiableRepository[T]) prepareVerifiableRecord(record T) error {
	firstRecord := false
	if strings.EqualFold(record.GetId(), "") {
		firstRecord = true
	}

	if !firstRecord {
		record.SetPrevious(record.GetId())
		record.SetSequenceNumber(record.GetSequenceNumber() + 1)
	}

	if r.noncer != nil {
		if err := record.GenerateNonce(r.noncer); err != nil {
			return err
		}
	}

	if r.timestamp {
		record.StampCreatedAt(nil)
	}

	if firstRecord {
		if err := algorithms.CreatePrefix(record); err != nil {
			return err
		}
	} else {
		if err := algorithms.SelfAddress(record); err != nil {
			return err
		}
	}

	return nil
}

func (r VerifiableRepository[T]) verifyRecord(record T) error {
	if record.GetSequenceNumber() == 0 {
		if err := algorithms.VerifyPrefixAndData(record); err != nil {
			return err
		}
	} else {
		if err := algorithms.VerifyAddressAndData(record); err != nil {
			return err
		}
	}

	return nil
}

// sql helpers

func (r VerifiableRepository[T]) insertRecord(ctx context.Context, record T) error {
	fieldNames := r.getFieldNames(record)
	innerFields := strings.Join(fieldNames, ", ")
	innerValues := strings.Join(fieldNames, ", :")

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

func (r VerifiableRepository[T]) getRecordById(ctx context.Context, record T, id string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %s", record.TableName(), r.store.Placeholders(1)[0])

	if err := r.store.Sql().GetContext(ctx, record, query, id); err != nil {
		return err
	}

	return nil
}

func (r VerifiableRepository[T]) getLatestRecordByPrefix(ctx context.Context, record T, prefix string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE prefix = %s ORDER BY sequence_number DESC LIMIT 1", record.TableName(), r.store.Placeholders(1)[0])

	if err := r.store.Sql().GetContext(ctx, record, query, prefix); err != nil {
		return err
	}

	return nil
}

func (r VerifiableRepository[T]) listRecordsByPrefix(ctx context.Context, records *[]T, prefix string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE prefix = %s ORDER BY sequence_number ASC", (*new(T)).TableName(), r.store.Placeholders(1)[0])

	if err := r.store.Sql().SelectContext(ctx, records, query, prefix); err != nil {
		return err
	}

	return nil
}

// sql helper helpers

func (r VerifiableRepository[T]) getFieldNames(s T) (fields []string) {
	v := reflect.ValueOf(s)
	t := v.Type()
	return r.getLeafFieldNamesWithValues(t, v)
}

func (r VerifiableRepository[T]) getLeafFieldNamesWithValues(t reflect.Type, v reflect.Value) []string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		if v.Kind() == reflect.Pointer && !v.IsNil() {
			v = v.Elem()
		}
	}

	if t.Kind() != reflect.Struct {
		return []string{}
	}

	var names []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldVal := v.Field(i)

		if fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeOf(primitives.Timestamp{}) {
			nested := r.getLeafFieldNamesWithValues(fieldType, fieldVal)
			names = append(names, nested...)
			continue
		}

		tag := field.Tag.Get("db")

		if strings.HasSuffix(tag, ",omitempty") && fieldVal.IsNil() {
			continue
		}
		if tag == "-" {
			continue
		}

		if tag == "" {
			names = append(names, field.Name)
		} else {
			names = append(names, strings.TrimSuffix(tag, ",omitempty"))
		}
	}
	return names
}
