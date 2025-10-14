package primitives

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
