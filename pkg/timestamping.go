package verifiablestorage

import (
	"encoding/json"
	"time"
)

const ConsistentNano = `2006-01-02T15:04:05.000000000Z07:00`

type Timestamp time.Time

func (t Timestamp) UTC() Timestamp {
	utc := (time.Time(t)).UTC()
	return Timestamp(utc)
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	when := time.Time(t).Format(ConsistentNano)
	b, err := json.Marshal(when)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	when, err := time.Parse(ConsistentNano, s)
	if err != nil {
		return err
	}

	*t = Timestamp(when)

	return nil
}

type Timestampable interface {
	// if when is null, Now() is used
	StampCreatedAt(when *Timestamp)
}

type Timestamper struct {
	CreatedAt Timestamp `db:"created_at" json:"createdAt"`
}

func (t *Timestamper) StampCreatedAt(when *Timestamp) {
	if when == nil {
		now := time.Now()
		when = (*Timestamp)(&now)
	}

	t.CreatedAt = when.UTC()
}
