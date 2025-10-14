package primitives

type Sequenceable interface {
	GetSequenceNumber() int
	SetSequenceNumber(sequenceNumber int)
}

type Sequencer struct {
	SequenceNumber int `db:"sequence_number" json:"sequenceNumber"`
}

func (s *Sequencer) IncrementSequenceNumber() {
	s.SequenceNumber += 1
}

func (s Sequencer) GetSequenceNumber() int {
	return s.SequenceNumber
}

func (s *Sequencer) SetSequenceNumber(sequenceNumber int) {
	s.SequenceNumber = sequenceNumber
}
