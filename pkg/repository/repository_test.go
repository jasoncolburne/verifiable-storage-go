package repository_test

import (
	"fmt"
	"testing"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces/examples"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/repository"
)

type SignableSearchableModel struct {
	primitives.SignableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
	Baz string `db:"baz" json:"baz"`
}

func (s SignableSearchableModel) DeriveSearchKey() string {
	// in the db, you could create a unique constraint around (search_key, sequence_number)
	blake3 := examples.NewBlake3()

	return blake3.Sum(s.Foo + s.Bar)
}

func TestSignableSearchableRepository(t *testing.T) {
	if err := exerciseSignableSearchableRepository(); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func exerciseSignableSearchableRepository() error {
	var store data.Store
	key, err := examples.NewEd25519(nil)
	if err != nil {
		return err
	}

	noncer := examples.NewNoncer()

	repository := repository.NewSignableRepository[*SignableSearchableModel](
		store,
		true,
		noncer,
		key,
	)

	record := &SignableSearchableModel{
		Foo: "bar",
		Bar: "baz",
		Baz: "foo",
	}

	repository.CreateVersion(record)
	repository.CreateVersion(record)
	repository.CreateVersion(record)

	if record.GetSequenceNumber().Int64() != 2 {
		return fmt.Errorf("unexpected sequence number")
	}

	return nil
}
