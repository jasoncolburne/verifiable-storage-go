package verifiablestorage

import (
	"encoding/json"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/crypto"
)

type Signable interface {
	SetSignature(signature string)
	GetSignature() string
	SetSigningIdentity(identity string)
	GetSigningIdentity() string
}

type Signer struct {
	SigningIdentity string `db:"signing_identity" json:"signingIdentity"`
	Signature       string `db:"signature" json:"-"`
}

func (s *Signer) SetSignature(signature string) {
	s.Signature = signature
}

func (s Signer) GetSignature() string {
	return s.Signature
}

func (s *Signer) SetSigningIdentity(identity string) {
	s.SigningIdentity = identity
}

func (s Signer) GetSigningIdentity() string {
	return s.SigningIdentity
}

func Sign(s Signable, key crypto.SigningKey) error {
	identity, err := key.Identity()
	if err != nil {
		return err
	}

	s.SetSigningIdentity(identity)

	message, err := json.Marshal(s)
	if err != nil {
		return err
	}

	signature, err := key.Sign(message)
	if err != nil {
		return err
	}

	s.SetSignature(signature)

	return nil
}

func VerifySignature(s Signable, verificationKeyStore crypto.VerificationKeyStore) error {
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
