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

	if !strings.EqualFold(string(jsonRecord), `{"id":"EMQ1elB_Sw3i73owirOMUO8z6Qagun3iFwOggMjeIfvF","prefix":"EMQ1elB_Sw3i73owirOMUO8z6Qagun3iFwOggMjeIfvF","sequenceNumber":0,"nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:32.722Z","foo":"bar","bar":"foo"}`) {
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

	if !strings.EqualFold(string(jsonRecord), `{"id":"EBV02CaOIRhCVwFv6EL3Ju6vb3ZQJY9m2F6BDsefrTuC","prefix":"EMQ1elB_Sw3i73owirOMUO8z6Qagun3iFwOggMjeIfvF","sequenceNumber":1,"previous":"EMQ1elB_Sw3i73owirOMUO8z6Qagun3iFwOggMjeIfvF","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:33.722Z","foo":"bar","bar":"foo"}`) {
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

	if !strings.EqualFold(string(jsonRecord), `{"id":"EPRTxjr3ABOKyMqkQQYvfoICjndmMpuY3CsEmUzSEgHU","prefix":"EMQ1elB_Sw3i73owirOMUO8z6Qagun3iFwOggMjeIfvF","sequenceNumber":2,"previous":"EBV02CaOIRhCVwFv6EL3Ju6vb3ZQJY9m2F6BDsefrTuC","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:34.722Z","foo":"bar","bar":"foo"}`) {
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

	if !strings.EqualFold(string(jsonRecord), `{"id":"ECs1y_ZbBx4P63gkuCm3dj8iJF1ILYNTW6c3EKhLWscU","prefix":"ECs1y_ZbBx4P63gkuCm3dj8iJF1ILYNTW6c3EKhLWscU","sequenceNumber":0,"nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:32.722Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 0 had unexpected values: %s", jsonRecord)
	}

	if !strings.EqualFold(r.Signature, "0BB9gq_PkM84YmBnK7IdHqeA7yLY_PYY5ulAPWGXjdq41zwW40sy8STD89Jr_9N8Ho4AzfPt_yB9xEN5pBJsa0QD") {
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

	if !strings.EqualFold(string(jsonRecord), `{"id":"EDWt5yTMe2qE9QyV0QIMSjOC9Nmvzw5OAYPBbM_m5GAV","prefix":"ECs1y_ZbBx4P63gkuCm3dj8iJF1ILYNTW6c3EKhLWscU","sequenceNumber":1,"previous":"ECs1y_ZbBx4P63gkuCm3dj8iJF1ILYNTW6c3EKhLWscU","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:33.722Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 1 had unexpected values: %s", jsonRecord)
	}

	if !strings.EqualFold(r.Signature, "0BDMUwphicKtgp4NSeFUsjy8EXYwPGjVKPRCxzTLxf2R2p10djQXAITBBQEFzC76-Q88VfWS0IauB54yis69TD8B") {
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

	if !strings.EqualFold(string(jsonRecord), `{"id":"EG2jlz7Lo6ISSKHI1pVJSY8sTMNDrWyKeKs3OldACTg7","prefix":"ECs1y_ZbBx4P63gkuCm3dj8iJF1ILYNTW6c3EKhLWscU","sequenceNumber":2,"previous":"EDWt5yTMe2qE9QyV0QIMSjOC9Nmvzw5OAYPBbM_m5GAV","nonce":"0A0000000000000000000000","createdAt":"2025-10-13T20:25:34.722Z","signingIdentity":"BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop","foo":"bar","bar":"foo"}`) {
		return fmt.Errorf("sn 2 had unexpected values: %s", jsonRecord)
	}

	if !strings.EqualFold(r.Signature, "0BACTCsjVbf6Bw-UKi-kbj3NKjc816HgFdLruGc62xRMA2XCEJpIv9LTBRFO_MU6eTDiROY_5-Hx4vPNCGk8rWIE") {
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
