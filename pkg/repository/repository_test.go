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

func (t SignableModel) TableName() string {
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

func TestSignableSearchableRepository(t *testing.T) {
	if err := exerciseSignableSearchableRepository(); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func exerciseSignableSearchableRepository() error {
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

	noncer := examples.NewNoncer()

	repository := repository.NewSignableRepository[*SignableModel](
		store,
		true,
		noncer,
		key,
	)

	record := &SignableModel{
		Foo: "bar",
		Bar: "baz",
	}

	if err := repository.CreateVersion(ctx, record); err != nil {
		return err
	}

	id := record.Id

	if err := repository.CreateVersion(ctx, record); err != nil {
		return err
	}

	if err := repository.CreateVersion(ctx, record); err != nil {
		return err
	}

	reloadedRecord := &SignableModel{}
	if err := repository.GetById(ctx, reloadedRecord, id); err != nil {
		return err
	}

	if !strings.EqualFold(reloadedRecord.Prefix, record.Prefix) {
		return fmt.Errorf("mismatched prefixes")
	}

	return nil
}
