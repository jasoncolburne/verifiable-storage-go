package primitives

type Sequenceable interface {
	GetSequenceNumber() uint64
	SetSequenceNumber(sequenceNumber uint64)
}

type Sequencer struct {
	SequenceNumber uint64 `db:"sequence_number" json:"sequenceNumber"`
}

func (s Sequencer) GetSequenceNumber() uint64 {
	return s.SequenceNumber
}

func (s *Sequencer) SetSequenceNumber(sequenceNumber uint64) {
	s.SequenceNumber = sequenceNumber
}
