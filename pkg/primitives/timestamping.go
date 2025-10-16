package primitives

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Timestamp time.Time

type Timestampable interface {
	// if when is null, Now() is used
	StampCreatedAt(when *Timestamp)
}

type Timestamper struct {
	CreatedAt *Timestamp `db:"created_at,omitempty" json:"createdAt,omitempty"`
}

func (t *Timestamper) StampCreatedAt(when *Timestamp) {
	if when == nil {
		now := time.Now()
		when = (*Timestamp)(&now)
	}

	utc := when.UTC()
	t.CreatedAt = &utc
}

const ConsistentNano = `2006-01-02T15:04:05.000000000Z07:00`

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

func (t Timestamp) Value() (driver.Value, error) {
	return time.Time(t).Format(ConsistentNano), nil
}

func (t *Timestamp) Scan(src any) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		*t = Timestamp(v)
		return nil
	case string:
		_t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}

		*t = Timestamp(_t)
		return nil
	case []byte:
		_t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}

		*t = Timestamp(_t)
		return nil
	default:
		return fmt.Errorf("unsupported src type %T", src)
	}
}
