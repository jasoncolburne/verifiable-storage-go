package crypto

type VerificationKeyStore interface {
	Get(identity string) (VerificationKey, error)
}
