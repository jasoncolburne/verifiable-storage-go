package interfaces

type Verifier interface {
	Verify(signature, publicKey string, message []byte) error
}
