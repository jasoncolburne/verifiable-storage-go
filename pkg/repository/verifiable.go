package repository

import (
	"math/big"
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

func (v VerifiableRepository[T]) CreateVersion(record T) error {
	v.prepareRecord(record, v.noncer)

	if v.write {
		// write to data store
	}

	return nil
}

func (v VerifiableRepository[T]) prepareRecord(record T, noncer interfaces.Noncer) error {
	firstRecord := false
	if strings.EqualFold(record.GetId(), "") {
		firstRecord = true
	}

	if !firstRecord {
		record.SetPrevious(record.GetId())
		sequenceNumber := record.GetSequenceNumber()
		sequenceNumber.Add(&sequenceNumber, big.NewInt(1))
		record.SetSequenceNumber(sequenceNumber)
	}

	if err := record.GenerateNonce(noncer); err != nil {
		return err
	}

	record.StampCreatedAt(nil)

	if firstRecord {
		if err := primitives.CreatePrefix(record); err != nil {
			return err
		}
	} else {
		if err := primitives.SelfAddress(record); err != nil {
			return err
		}
	}

	return nil
}
