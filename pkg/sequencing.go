package verifiablestorage

type Sequenceable interface {
	IncrementSequenceNumber()
}

type Sequencer struct {
	SequenceNumber int `db:"sequence_number" json:"sequenceNumber"`
}

func (s *Sequencer) IncrementSequenceNumber() {
	s.SequenceNumber += 1
}
