package algorithms_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces/examples"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

func TestSigning(t *testing.T) {
	if err := testSigning(); err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
}

func testSigning() error {
	signer := &primitives.Signer{}

	seed := [32]byte{}

	key, err := examples.NewEd25519(seed[:])
	if err != nil {
		return err
	}

	identity, err := key.Identity()
	if err != nil {
		return err
	}

	keyStore := examples.NewVerificationKeyStore()
	keyStore.Add(identity, key)

	if err := algorithms.Sign(signer, key, func() error { return nil }); err != nil {
		return err
	}

	expectedIdentity := `BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdop`
	badIdentity := `BDtqJ7zOtqQtYqOo0CpvDXNlMhV3HeJDpjrASKGLWdoq`

	expectedSignature := `0BC_U_WXfHTqyZL-nFhHvCi5w3oEjkWwSPOtlKBF-9uIQrZUU6qlwjdubPah-Hw7mJB7pbLbZNNEk5scR_uPCZ0K`
	badSignature := `0BC_U_WXfHTqyZL-nFhHvCi5w3oEjkWwSPOtlKBF-9uIQrZUU6qlwjdubPah-Hw7mJB7pbLbZNNEk5scR_uPCZ0L`

	if !strings.EqualFold(signer.SigningIdentity, expectedIdentity) {
		return fmt.Errorf("unexpected identity: %s", signer.SigningIdentity)
	}

	if !strings.EqualFold(signer.Signature, expectedSignature) {
		return fmt.Errorf("unexpected signature: %s", signer.Signature)
	}

	if err := algorithms.VerifySignature(signer, keyStore); err != nil {
		return err
	}

	signer.Signature = badSignature

	if err := algorithms.VerifySignature(signer, keyStore); err == nil {
		return fmt.Errorf("unexpected verification success for bad signature")
	}

	signer.Signature = expectedSignature
	signer.SigningIdentity = badIdentity

	if err := algorithms.VerifySignature(signer, keyStore); err == nil {
		return fmt.Errorf("unexpected verification success for bad identity")
	}

	return nil
}
