package interfaces

type VerificationKey interface {
	Verifier() Verifier
	Public() (string, error)
}
