package repository_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	data "github.com/jasoncolburne/verifiable-storage-go/pkg/data/examples"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces/examples"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/repository"
)

type VerifiableModel struct {
	primitives.VerifiableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
}

func (*VerifiableModel) TableName() string {
	return `verifiable`
}

var VERIFIABLE_TABLE_SQL = `
CREATE TABLE IF NOT EXISTS verifiable (
	-- Standard fields
    id              	TEXT PRIMARY KEY,
	prefix				TEXT NOT NULL,
	previous        	TEXT,
	sequence_number 	BIGINT NOT NULL,
	created_at          DATETIME NOT NULL,
	nonce           	TEXT NOT NULL,

	-- Model-specific fields
	foo 				TEXT NOT NULL,
	bar                 TEXT NOT NULL,

	-- Uniqueness constraint for sequence numbers
	UNIQUE(prefix, sequence_number)
);
`

type SignableModel struct {
	primitives.SignableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
}

func (*SignableModel) TableName() string {
	return `signable`
}

var SIGNABLE_TABLE_SQL = `
CREATE TABLE IF NOT EXISTS signable (
	-- Standard fields
    id              	TEXT PRIMARY KEY,
	prefix				TEXT NOT NULL,
	previous        	TEXT,
	sequence_number 	BIGINT NOT NULL,
	created_at          DATETIME NOT NULL,
	nonce           	TEXT NOT NULL,

	-- Signable fields
	signing_identity	TEXT NOT NULL,
	signature       	TEXT NOT NULL,

	-- Model-specific fields
	foo 				TEXT NOT NULL,
	bar                 TEXT NOT NULL,

	-- Uniqueness constraint for sequence numbers
	UNIQUE(prefix, sequence_number)
);
`

func TestVerifiableRepository(t *testing.T) {
	repository, err := createVerifiableRepository()
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}

	record := &VerifiableModel{
		Foo: "bar",
		Bar: "baz",
	}

	buffers := []*VerifiableModel{{}, {}, {}}

	if err := exerciseRepository(repository, record, buffers); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func createVerifiableRepository() (repository.Repository[*VerifiableModel], error) {
	ctx := context.Background()

	store, err := data.NewInMemorySQLiteStore()
	if err != nil {
		return nil, err
	}

	_, err = store.Sql().ExecContext(ctx, VERIFIABLE_TABLE_SQL)
	if err != nil {
		return nil, err
	}

	noncer := examples.NewNoncer()

	repository := repository.NewVerifiableRepository[*VerifiableModel](
		store,
		true,
		noncer,
	)

	return repository, nil
}

func TestSignableRepository(t *testing.T) {
	repository, err := createSignableRepository()
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}

	record := &SignableModel{
		Foo: "bar",
		Bar: "baz",
	}

	buffers := []*SignableModel{{}, {}, {}}

	if err := exerciseRepository(repository, record, buffers); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func createSignableRepository() (repository.Repository[*SignableModel], error) {
	ctx := context.Background()

	store, err := data.NewInMemorySQLiteStore()
	if err != nil {
		return nil, err
	}

	_, err = store.Sql().ExecContext(ctx, SIGNABLE_TABLE_SQL)
	if err != nil {
		return nil, err
	}

	key, err := examples.NewEd25519(nil)
	if err != nil {
		return nil, err
	}

	keyIdentity, err := key.Identity()
	if err != nil {
		return nil, err
	}

	verificationKeyStore := examples.NewVerificationKeyStore()
	verificationKeyStore.Add(keyIdentity, key)

	noncer := examples.NewNoncer()

	repository := repository.NewSignableRepository[*SignableModel](
		store,
		true,
		noncer,
		key,
		verificationKeyStore,
	)

	return repository, nil
}

func exerciseRepository[T primitives.VerifiableAndRecordable](repository repository.Repository[T], record2 T, buffers []T) error {
	ctx := context.Background()

	if err := repository.CreateVersion(ctx, record2); err != nil {
		return err
	}

	if err := repository.CreateVersion(ctx, record2); err != nil {
		return err
	}

	id1 := record2.GetId()

	if err := repository.CreateVersion(ctx, record2); err != nil {
		return err
	}

	record1 := buffers[0]
	if err := repository.GetById(ctx, record1, id1); err != nil {
		return err
	}

	record0 := buffers[1]
	if err := repository.GetById(ctx, record0, record2.GetPrefix()); err != nil {
		return err
	}

	latest := buffers[2]
	if err := repository.GetLatestByPrefix(ctx, latest, record0.GetId()); err != nil {
		return err
	}

	if record0.GetSequenceNumber() != 0 {
		return fmt.Errorf("unexpected sn for 0: %d", record0.GetSequenceNumber())
	}

	if record0.GetPrevious() != nil {
		return fmt.Errorf("previous not nil")
	}

	if !strings.EqualFold(record1.GetPrefix(), record2.GetPrefix()) {
		return fmt.Errorf("mismatched prefixes")
	}

	if strings.EqualFold(record1.GetId(), record2.GetId()) {
		return fmt.Errorf("unexpected equal ids")
	}

	if record1.GetSequenceNumber() != 1 {
		return fmt.Errorf("unexpected sn for 1: %d", record1.GetSequenceNumber())
	}

	if record2.GetSequenceNumber() != 2 {
		return fmt.Errorf("unexpected sn for 2: %d", record2.GetSequenceNumber())
	}

	if record2.GetPrevious() == nil || !strings.EqualFold(*record2.GetPrevious(), record1.GetId()) {
		return fmt.Errorf("mismatched previous 1")
	}

	if record1.GetPrevious() == nil || !strings.EqualFold(*record1.GetPrevious(), record1.GetPrefix()) {
		return fmt.Errorf("mismatched previous 0")
	}

	if !strings.EqualFold(latest.GetId(), record2.GetId()) {
		return fmt.Errorf("latest id is mismatched")
	}

	records := []T{}

	if err := repository.ListByPrefix(ctx, &records, record2.GetPrefix()); err != nil {
		return err
	}

	if len(records) != 3 {
		return fmt.Errorf("incorrect number of records listed (%d)", len(records))
	}

	for i, record := range records {
		switch i {
		case 0:
			if !strings.EqualFold(record0.GetId(), record.GetId()) {
				return fmt.Errorf("listed record 0 has incorrect id")
			}
		case 1:
			if !strings.EqualFold(record1.GetId(), record.GetId()) {
				return fmt.Errorf("listed record 1 has incorrect id")
			}
		case 2:
			if !strings.EqualFold(record2.GetId(), record.GetId()) {
				return fmt.Errorf("listed record 2 has incorrect id")
			}
		}
	}

	return nil
}
