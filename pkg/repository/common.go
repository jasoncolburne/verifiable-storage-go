package repository

import (
	"math/big"
	"strings"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

func prepareVerifiableRecord(record primitives.VerifiableAndRecordable, noncer interfaces.Noncer) error {
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

func prepareSearchableRecord(record primitives.SearchableAndRecordable, noncer interfaces.Noncer) error {
	record.SetSearchKey(record.DeriveSearchKey())

	return prepareVerifiableRecord(record, noncer)
}

func prepareSignableRecord(record primitives.SignableAndRecordable, noncer interfaces.Noncer, key interfaces.SigningKey) error {
	if err := algorithms.Sign(record, key); err != nil {
		return err
	}

	return prepareVerifiableRecord(record, noncer)
}

func prepareSignableSearchableRecord(record primitives.SignableAndSearchableAndRecordable, noncer interfaces.Noncer, key interfaces.SigningKey) error {
	record.SetSearchKey(record.DeriveSearchKey())

	return prepareSignableRecord(record, noncer, key)
}
