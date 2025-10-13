package crypto

type Verifier interface {
	Verify(signature, publicKey string, message []byte) error
}
