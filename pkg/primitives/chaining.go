package primitives

type Chainable interface {
	GetPrevious() *string
	SetPrevious(previous string)
}

type Chainer struct {
	Previous *string `db:"previous" json:"previous,omitempty"`
}

func (c Chainer) GetPrevious() *string {
	return c.Previous
}

func (c *Chainer) SetPrevious(previous string) {
	c.Previous = &previous
}
