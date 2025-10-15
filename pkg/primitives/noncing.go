package primitives

import "github.com/jasoncolburne/verifiable-storage-go/pkg/interfaces"

type Nonceable interface {
	GenerateNonce(source interfaces.Noncer) error
}

type Noncer struct {
	Nonce string `db:"nonce" json:"nonce"`
}

func (n *Noncer) GenerateNonce(source interfaces.Noncer) error {
	nonce, err := source.Generate()
	if err != nil {
		return err
	}

	n.Nonce = nonce

	return nil
}
