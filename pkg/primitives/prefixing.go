package primitives

type Prefixable interface {
	SelfAddressable
	GetPrefix() string
	SetPrefix(prefix string)
}

type Prefixer struct {
	SelfAddresser
	Prefix string `db:"prefix" json:"prefix"`
}

func (p Prefixer) GetPrefix() string {
	return p.Prefix
}

func (p *Prefixer) SetPrefix(prefix string) {
	p.Prefix = prefix
}
