package verifiablestorage

import (
	"encoding/json"
	"fmt"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/crypto"
)

type Signable interface {
	Sign(key crypto.SigningKey) error
	Verify(verificationKeyStore crypto.VerificationKeyStore) error
}

type Signer struct {
	SigningIdentity string `db:"signing_identity" json:"signingIdentity"`
	Signature       string `db:"signature" json:"-"`
}

func (s *Signer) Sign(key crypto.SigningKey) error {
	message, err := json.Marshal(s)
	if err != nil {
		return err
	}

	s.SigningIdentity, err = key.Identity()
	if err != nil {
		return err
	}

	s.Signature, err = key.Sign(message)
	return err
}

func (s Signer) Verify(verificationKeyStore crypto.VerificationKeyStore) error {
	verificationKey, err := verificationKeyStore.Get(s.SigningIdentity)
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

	fmt.Printf("%s\n", message)

	return verificationKey.Verifier().Verify(s.Signature, verificationPublicKey, message)
}
