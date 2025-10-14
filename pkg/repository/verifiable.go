package repository

import (
	"encoding/json"
	"fmt"

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

func (r VerifiableRepository[T]) CreateVersion(record T) error {
	prepareVerifiableRecord(record, r.noncer)

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
