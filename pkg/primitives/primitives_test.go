package primitives_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces/examples"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type Record struct {
	primitives.VerifiableRecorder
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

	if err := createVersion(r, noncer); err != nil {
		return err
	}

	if err := createVersion(r, noncer); err != nil {
		return err
	}

	if r.SequenceNumber != 2 {
		return fmt.Errorf("not incremented")
	}

	return nil
}

func createVersion(r primitives.VerifiableAndRecordable, noncer interfaces.Noncer) error {
	firstRecord := false
	if strings.EqualFold(r.GetId(), "") {
		firstRecord = true
	}

	if !firstRecord {
		r.SetPrevious(r.GetId())
		r.IncrementSequenceNumber()
	}

	if err := r.GenerateNonce(noncer); err != nil {
		return err
	}

	r.StampCreatedAt(nil)

	if firstRecord {
		if err := primitives.CreatePrefix(r); err != nil {
			return err
		}
	} else {
		if err := primitives.SelfAddress(r); err != nil {
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

	if err := createFixedVerifiableVersion(r, "2025-10-13T20:25:32.722276000Z"); err != nil {
		return err
	}

	jsonRecord, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","prefix":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","sequenceNumber":0,"nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:32.722276000Z","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 0 had unexpected values")
	}

	if err := createFixedVerifiableVersion(r, "2025-10-13T20:25:33.722276000Z"); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EO_wM1UoWjk4aTOkrOdN56JxfNJBwGpKDpFAcaSlEiB3","prefix":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","sequenceNumber":1,"previous":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:33.722276000Z","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 1 had unexpected values")
	}

	if err := createFixedVerifiableVersion(r, "2025-10-13T20:25:34.722276000Z"); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EK-eJ0YStKtbjNoLeFUrC1secwP9rtqE4gQ9_cKKwmuu","prefix":"EKV6bgU7UuFzQYqsvovO2yPV6r6pZss6OzxpJJgI0HEq","sequenceNumber":2,"previous":"EO_wM1UoWjk4aTOkrOdN56JxfNJBwGpKDpFAcaSlEiB3","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:34.722276000Z","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 2 had unexpected values")
	}

	return nil
}

func createFixedVerifiableVersion(r primitives.VerifiableAndRecordable, at string) error {
	noncer := &FixedNoncer{}

	firstRecord := false
	if strings.EqualFold(r.GetId(), "") {
		firstRecord = true
	}

	if !firstRecord {
		r.SetPrevious(r.GetId())
		r.IncrementSequenceNumber()
	}

	if err := r.GenerateNonce(noncer); err != nil {
		return err
	}

	when, err := time.Parse(time.RFC3339Nano, at)
	if err != nil {
		return err
	}
	r.StampCreatedAt((*primitives.Timestamp)(&when))

	if firstRecord {
		if err := primitives.CreatePrefix(r); err != nil {
			return err
		}
	} else {
		if err := primitives.SelfAddress(r); err != nil {
			return err
		}
	}

	return nil
}

type SignableRecord struct {
	primitives.SignableRecorder
	Foo string `db:"foo" json:"foo"`
	Bar string `db:"bar" json:"bar"`
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

	if err := createFixedSignedVersion(r, "2025-10-13T20:25:32.722276000Z", key); err != nil {
		return err
	}

	if err := verifySignature(r, key); err != nil {
		return err
	}

	jsonRecord, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"ELG3AqCIt2FgklHK_TI3dXVLqlHlxb9v2Kvl-IQ4Hhgo","prefix":"ELG3AqCIt2FgklHK_TI3dXVLqlHlxb9v2Kvl-IQ4Hhgo","sequenceNumber":0,"nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:32.722276000Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 0 had unexpected values")
	}

	if !strings.EqualFold(r.Signature, "0BDgyCjwCxzSBWx-SuPez_VXIZbWNW8wwBzFG4tGVD1jqG1HidhLYXo6lGzSC4gKzgVjif64wAUhFUoUdgfZRs0D") {
		return fmt.Errorf("sn 0 had unexpected signature: %s", r.Signature)
	}

	if err := createFixedSignedVersion(r, "2025-10-13T20:25:33.722276000Z", key); err != nil {
		return err
	}

	if err := verifySignature(r, key); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"EMmbYhG3GtJI52WTC06Z6s9gkIIVMhhQKP1-fotvuyDP","prefix":"ELG3AqCIt2FgklHK_TI3dXVLqlHlxb9v2Kvl-IQ4Hhgo","sequenceNumber":1,"previous":"ELG3AqCIt2FgklHK_TI3dXVLqlHlxb9v2Kvl-IQ4Hhgo","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:33.722276000Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 1 had unexpected values")
	}

	if !strings.EqualFold(r.Signature, "0BD0X3oVgKWf1F9N2l8AFcu4lKRdy5J895m21ytoQ8XTi4BxDSEo09gDYrU3owP0xgiSzXYnFrwHACVZwepWR3AG") {
		return fmt.Errorf("sn 1 had unexpected signature: %s", r.Signature)
	}

	if err := createFixedSignedVersion(r, "2025-10-13T20:25:34.722276000Z", key); err != nil {
		return err
	}

	if err := verifySignature(r, key); err != nil {
		return err
	}

	jsonRecord, err = json.Marshal(r)
	if err != nil {
		return err
	}

	if !strings.EqualFold(string(jsonRecord), `{"id":"ECjfjyfLwO3Cip1Q950zd2MjXZVokjG1mkdNKTL4rNO_","prefix":"ELG3AqCIt2FgklHK_TI3dXVLqlHlxb9v2Kvl-IQ4Hhgo","sequenceNumber":2,"previous":"EMmbYhG3GtJI52WTC06Z6s9gkIIVMhhQKP1-fotvuyDP","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:34.722276000Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		fmt.Printf("%s\n", jsonRecord)
		return fmt.Errorf("sn 2 had unexpected values")
	}

	if !strings.EqualFold(r.Signature, "0BC9OYswxd7AMxZPi6kGteFo48_At0KlivhcLJmzSP37Jc2n5yWqlYvz1yOQrM8UPCheryNincuaA2ms5Vpxi1MJ") {
		return fmt.Errorf("sn 2 had unexpected signature: %s", r.Signature)
	}

	return nil
}

func createFixedSignedVersion(s primitives.SignableAndRecordable, at string, key interfaces.SigningKey) error {
	if err := createFixedVerifiableVersion(s, at); err != nil {
		return err
	}

	if err := primitives.Sign(s, key); err != nil {
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

	if err := primitives.VerifySignature(s, verificationKeyStore); err != nil {
		return err
	}

	return nil
}
