package crypto

type VerificationKey interface {
	Verifier() Verifier
	Public() (string, error)
}
