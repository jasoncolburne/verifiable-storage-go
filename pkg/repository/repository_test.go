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

func TestSignableRepository(t *testing.T) {
	if err := exerciseSignableRepository(); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func exerciseSignableRepository() error {
	ctx := context.Background()

	store, err := data.NewInMemorySQLiteStore()
	if err != nil {
		return err
	}

	_, err = store.Sql().ExecContext(ctx, SIGNABLE_TABLE_SQL)
	if err != nil {
		return err
	}

	key, err := examples.NewEd25519(nil)
	if err != nil {
		return err
	}

	keyIdentity, err := key.Identity()
	if err != nil {
		return err
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

	record2 := &SignableModel{
		Foo: "bar",
		Bar: "baz",
	}

	if err := repository.CreateVersion(ctx, record2); err != nil {
		return err
	}

	if err := repository.CreateVersion(ctx, record2); err != nil {
		return err
	}

	id1 := record2.Id

	if err := repository.CreateVersion(ctx, record2); err != nil {
		return err
	}

	record1 := &SignableModel{}
	if err := repository.GetById(ctx, record1, id1); err != nil {
		return err
	}

	record0 := &SignableModel{}
	if err := repository.GetById(ctx, record0, record2.Prefix); err != nil {
		return err
	}

	latest := &SignableModel{}
	if err := repository.GetLatestByPrefix(ctx, latest, record0.Id); err != nil {
		return err
	}

	if record0.SequenceNumber != 0 {
		return fmt.Errorf("unexpected sn for 0: %d", record0.SequenceNumber)
	}

	if record0.Previous != nil {
		return fmt.Errorf("previous not nil")
	}

	if !strings.EqualFold(record1.Prefix, record2.Prefix) {
		return fmt.Errorf("mismatched prefixes")
	}

	if strings.EqualFold(record1.Id, record2.Id) {
		return fmt.Errorf("unexpected equal ids")
	}

	if record1.SequenceNumber != 1 {
		return fmt.Errorf("unexpected sn for 1: %d", record1.SequenceNumber)
	}

	if record2.SequenceNumber != 2 {
		return fmt.Errorf("unexpected sn for 2: %d", record2.SequenceNumber)
	}

	if record2.Previous == nil || !strings.EqualFold(*record2.Previous, record1.Id) {
		return fmt.Errorf("mismatched previous 1")
	}

	if record1.Previous == nil || !strings.EqualFold(*record1.Previous, record1.Prefix) {
		return fmt.Errorf("mismatched previous 0")
	}

	if !strings.EqualFold(latest.Id, record2.Id) {
		return fmt.Errorf("latest id is mismatched")
	}

	records := []*SignableModel{}

	if err := repository.ListByPrefix(ctx, &records, record2.Prefix); err != nil {
		return err
	}

	if len(records) != 3 {
		return fmt.Errorf("incorrect number of records listed (%d)", len(records))
	}

	for i, record := range records {
		switch i {
		case 0:
			if !strings.EqualFold(record0.Id, record.Id) {
				return fmt.Errorf("listed record 0 has incorrect id")
			}
		case 1:
			if !strings.EqualFold(record1.Id, record.Id) {
				return fmt.Errorf("listed record 1 has incorrect id")
			}
		case 2:
			if !strings.EqualFold(record2.Id, record.Id) {
				return fmt.Errorf("listed record 2 has incorrect id")
			}
		}
	}

	return nil
}
