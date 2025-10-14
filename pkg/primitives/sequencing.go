package primitives

import "math/big"

type Sequenceable interface {
	GetSequenceNumber() *big.Int
	SetSequenceNumber(sequenceNumber *big.Int)
}

type Sequencer struct {
	SequenceNumber big.Int `db:"sequence_number" json:"sequenceNumber"`
}

func (s Sequencer) GetSequenceNumber() *big.Int {
	return &s.SequenceNumber
}

func (s *Sequencer) SetSequenceNumber(sequenceNumber *big.Int) {
	s.SequenceNumber = *sequenceNumber
}
