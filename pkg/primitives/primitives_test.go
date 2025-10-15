package primitives_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces/examples"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type Record struct {
	primitives.VerifiableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
}

func (Record) TableName() string {
	return `record`
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

	if err := createVerifiableVersion(r, nil, noncer); err != nil {
		return err
	}

	if err := createVerifiableVersion(r, nil, noncer); err != nil {
		return err
	}

	if err := createVerifiableVersion(r, nil, noncer); err != nil {
		return err
	}

	if r.SequenceNumber != 2 {
		return fmt.Errorf("not incremented")
	}

	return nil
}

func createVerifiableVersion(r primitives.VerifiableAndRecordable, at *primitives.Timestamp, noncer interfaces.Noncer) error {
	firstRecord := false
	if strings.EqualFold(r.GetId(), "") {
		firstRecord = true
	}

	if !firstRecord {
		r.SetPrevious(r.GetId())
		r.SetSequenceNumber(r.GetSequenceNumber() + 1)
	}

	if err := r.GenerateNonce(noncer); err != nil {
		return err
	}

	r.StampCreatedAt(at)

	if firstRecord {
		if err := algorithms.CreatePrefix(r); err != nil {
			return err
		}
	} else {
		if err := algorithms.SelfAddress(r); err != nil {
			return err
		}
	}

	return nil
}

type FixedNoncer struct{}

func (FixedNoncer) Generate() (string, error) {
	return "0A0000000000000000000000", nil
}

func TestFixedVerifiableRecorder(t *testing.T) {
	if err := exerciseFixedVerifiableRecorder(); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func exerciseFixedVerifiableRecorder() error {
	r := &Record{
		Foo: "bar",
		Bar: "foo",
	}

	at, err := time.Parse(time.RFC3339Nano, "2025-10-13T20:25:32.722276000Z")
	if err != nil {
		return err
	}

	if err := createFixedVerifiableVersion(r, primitives.Timestamp(at)); err != nil {
		return err
	}

	jsonRecord, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","prefix":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","sequenceNumber":0,"nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:32.722276000Z","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 0 had unexpected values: %s", jsonRecord)
	}

	at = at.Add(time.Second)
	if err := createFixedVerifiableVersion(r, primitives.Timestamp(at)); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EO_wM1UoWjk4aTOkrOdN56JxfNJBwGpKDpFAcaSlEiB3","prefix":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","sequenceNumber":1,"previous":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:33.722276000Z","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 1 had unexpected values: %s", jsonRecord)
	}

	at = at.Add(time.Second)
	if err := createFixedVerifiableVersion(r, primitives.Timestamp(at)); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EK-eJ0YStKtbjNoLeFUrC1secwP9rtqE4gQ9_cKKwmuu","prefix":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","sequenceNumber":2,"previous":"EO_wM1UoWjk4aTOkrOdN56JxfNJBwGpKDpFAcaSlEiB3","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:34.722276000Z","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 2 had unexpected values: %s", jsonRecord)
	}

	return nil
}

func createFixedVerifiableVersion(r primitives.VerifiableAndRecordable, at primitives.Timestamp) error {
	noncer := &FixedNoncer{}
	when := &at

	return createVerifiableVersion(r, when, noncer)
}

type SignableRecord struct {
	primitives.SignableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
}

func (SignableRecord) TableName() string {
	return `signablerecord`
}

func TestFixedSignableRecorder(t *testing.T) {
	if err := exerciseFixedSignableRecorder(); err != nil {
		fmt.Printf("%s\n", err)
		t.Fail()
	}
}

func exerciseFixedSignableRecorder() error {
	r := &SignableRecord{
		Foo: "bar",
		Bar: "foo",
	}

	seed := [32]byte{}

	key, err := examples.NewEd25519(seed[:])
	if err != nil {
		return err
	}

	at, err := time.Parse(time.RFC3339Nano, "2025-10-13T20:25:32.722276000Z")
	if err != nil {
		return err
	}

	if err := createFixedSignedVersion(r, primitives.Timestamp(at), key); err != nil {
		return err
	}

	if err := verifySignature(r, key); err != nil {
		return err
	}

	jsonRecord, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EN98smTUR2b4_XDoX7kSg6PhZa38X4tYS7rS_QfyrdfV","prefix":"EN98smTUR2b4_XDoX7kSg6PhZa38X4tYS7rS_QfyrdfV","sequenceNumber":0,"nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:32.722276000Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 0 had unexpected values: %s", jsonRecord)
	}

	if !strings.EqualFold(r.Signature, "0BDYzrwL2kLyV9-cqlUpZPGb3Knv41fi5XS5cmvFRJ0H69EijqO0h-gxfsuPbpurdO1yFHaTZgfurKN-ORjXZqEH") {
		return fmt.Errorf("sn 0 had unexpected signature: %s", r.Signature)
	}

	at = at.Add(time.Second)
	if err := createFixedSignedVersion(r, primitives.Timestamp(at), key); err != nil {
		return err
	}

	if err := verifySignature(r, key); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EHlunc8O09d5S-rHTD7IOHZ9Uo1vSfr1C6VAuEGJLnak","prefix":"EN98smTUR2b4_XDoX7kSg6PhZa38X4tYS7rS_QfyrdfV","sequenceNumber":1,"previous":"EN98smTUR2b4_XDoX7kSg6PhZa38X4tYS7rS_QfyrdfV","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:33.722276000Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 1 had unexpected values: %s", jsonRecord)
	}

	if !strings.EqualFold(r.Signature, "0BCRqU5w6jV4xSf__vgxwlxiAUTEwZoczLxHPx7IAqfTaY88ua6p5SB_4kZYBy2vKe-PnQXNcSteCuHQVB51LdcD") {
		return fmt.Errorf("sn 1 had unexpected signature: %s", r.Signature)
	}

	at = at.Add(time.Second)
	if err := createFixedSignedVersion(r, primitives.Timestamp(at), key); err != nil {
		return err
	}

	if err := verifySignature(r, key); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EETf1YcE2NeQFY9hFCCCjThhelLxUGSyLQ9e3BCCaSYd","prefix":"EN98smTUR2b4_XDoX7kSg6PhZa38X4tYS7rS_QfyrdfV","sequenceNumber":2,"previous":"EHlunc8O09d5S-rHTD7IOHZ9Uo1vSfr1C6VAuEGJLnak","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:34.722276000Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 2 had unexpected values: %s", jsonRecord)
	}

	if !strings.EqualFold(r.Signature, "0BBzjQIE3QOU9c81N6V8WuJ5W8gQoWh9DXTMccQKx1_EJfuJvVzqeMOY5E6l2_n75Ywy4Wx2e01f2TxLM4zQox0N") {
		return fmt.Errorf("sn 2 had unexpected signature: %s", r.Signature)
	}

	return nil
}

func createFixedSignedVersion(s primitives.SignableAndRecordable, at primitives.Timestamp, key interfaces.SigningKey) error {
	if err := algorithms.Sign(s, key, func() error {
		return createFixedVerifiableVersion(s, at)
	}); err != nil {
		return err
	}

	return nil
}

func verifySignature(s primitives.SignableAndRecordable, key interfaces.SigningKey) error {
	identity, err := key.Identity()
	if err != nil {
		return err
	}

	publicKey, err := key.Public()
	if err != nil {
		return err
	}

	verificationKey := examples.NewEd25519VerificationKey(publicKey)

	verificationKeyStore := examples.NewVerificationKeyStore()
	verificationKeyStore.Add(identity, verificationKey)

	if err := algorithms.VerifySignature(s, verificationKeyStore); err != nil {
		return err
	}

	return nil
}
