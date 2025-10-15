package repository

import (
	"reflect"
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
		record.SetSequenceNumber(record.GetSequenceNumber() + 1)
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

func prepareSignedRecord(record primitives.SignableAndRecordable, noncer interfaces.Noncer, key interfaces.SigningKey) error {
	if err := algorithms.Sign(record, key, func() error {
		return prepareVerifiableRecord(record, noncer)
	}); err != nil {
		return err
	}

	return nil
}

func verifyRecord(record primitives.VerifiableAndRecordable) error {
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

func verifySignedRecord(record primitives.SignableAndRecordable, verificationKeyStore interfaces.VerificationKeyStore) error {
	if err := algorithms.VerifySignature(record, verificationKeyStore); err != nil {
		return err
	}

	if err := verifyRecord(record); err != nil {
		return err
	}

	return nil
}

func getFieldNames(s any) (fields []string) {
	t := reflect.TypeOf(s)
	return getLeafFieldNames(t)
}

func getLeafFieldNames(t reflect.Type) (names []string) {
	// If pointer, deref.
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// Process only struct types.
	if t.Kind() != reflect.Struct {
		return names
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		// If field is itself a struct (and not a time.Time or slice/map), recurse.
		if fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeOf(primitives.Timestamp{}) {
			nested := getLeafFieldNames(fieldType)
			names = append(names, nested...)
		} else {
			// Use db tag if present.
			tag := field.Tag.Get("db")
			if tag != "" && tag != "-" {
				names = append(names, tag)
			} else {
				names = append(names, field.Name)
			}
		}
	}
	return names
}
