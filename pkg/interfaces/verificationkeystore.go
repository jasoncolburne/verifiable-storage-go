package interfaces

type VerificationKeyStore interface {
	Get(identity string) (VerificationKey, error)
}
