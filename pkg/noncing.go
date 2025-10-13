package verifiablestorage

import "github.com/jasoncolburne/verifiable-storage-go/pkg/crypto"

type Nonceable interface {
	GenerateNonce(source crypto.Noncer) error
}

type Noncer struct {
	Nonce string `db:"nonce" json:"nonce"`
}

func (n *Noncer) GenerateNonce(source crypto.Noncer) error {
	nonce, err := source.Generate()
	if err != nil {
		return err
	}

	n.Nonce = nonce

	return nil
}
