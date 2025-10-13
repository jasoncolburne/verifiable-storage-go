package crypto

type SigningKey interface {
	VerificationKey
	Identity() (string, error)
	Sign(message []byte) (string, error)
}
