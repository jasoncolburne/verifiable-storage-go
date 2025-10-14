package algorithms

import (
	"encoding/json"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

func Sign(s primitives.Signable, key interfaces.SigningKey) error {
	identity, err := key.Identity()
	if err != nil {
		return err
	}

	message, err := json.Marshal(s)
	if err != nil {
		return err
	}

	signature, err := key.Sign(message)
	if err != nil {
		return err
	}

	s.SetSigningIdentity(identity)
	s.SetSignature(signature)

	return nil
}

func VerifySignature(s primitives.Signable, verificationKeyStore interfaces.VerificationKeyStore) error {
	verificationKey, err := verificationKeyStore.Get(s.GetSigningIdentity())
	if err != nil {
		return err
	}

	verificationPublicKey, err := verificationKey.Public()
	if err != nil {
		return err
	}

	message, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return verificationKey.Verifier().Verify(s.GetSignature(), verificationPublicKey, message)
}

func CreateSignedContainer[T primitives.Signable](record T) (string, error) {
	container := primitives.SignedContainer[T]{
		Record:    record,
		Signature: record.GetSignature(),
	}

	jsonString, err := json.Marshal(container)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}
