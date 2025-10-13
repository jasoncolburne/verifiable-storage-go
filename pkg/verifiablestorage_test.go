package verifiablestorage_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	verifiablestorage "github.com/jasoncolburne/verifiable-storage-go/pkg"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/crypto"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/crypto/examples"
)

type Record struct {
	verifiablestorage.VerifiableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
}

func TestVerifiableRecorder(t *testing.T) {
	if err := exerciseVerifiableRecorder(); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func exerciseVerifiableRecorder() error {
	r := &Record{
		Foo: "bar",
		Bar: "foo",
	}

	noncer := examples.NewNoncer()

	if err := createVersion(r, noncer); err != nil {
		return err
	}

	jsonRecord, err := json.Marshal(r)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonRecord)

	if err := createVersion(r, noncer); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonRecord)

	if err := createVersion(r, noncer); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonRecord)

	return fmt.Errorf("printing...")
}

func createVersion(r verifiablestorage.VerifiableAndRecordable, noncer crypto.Noncer) error {
	firstRecord := false
	if strings.EqualFold(r.GetId(), "") {
		firstRecord = true
	}

	if !firstRecord {
		r.SetPrevious(r.GetId())
		r.IncrementSequenceNumber()
	}

	if firstRecord {
		r.SetPrefix("############################################")
	}

	r.StampCreatedAt(nil)
	if err := r.GenerateNonce(noncer); err != nil {
		return err
	}

	if err := verifiablestorage.SelfAddress(r); err != nil {
		return err
	}

	if firstRecord {
		r.SetPrefix(r.GetId())
	}

	return nil
}
