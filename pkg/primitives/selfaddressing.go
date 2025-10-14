package primitives

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zeebo/blake3"
)

type SelfAddressable interface {
	GetId() string
	SetId(id string)
}

type SelfAddresser struct {
	Id string `db:"id" json:"id"`
}

func (s SelfAddresser) GetId() string {
	return s.Id
}

func (s *SelfAddresser) SetId(id string) {
	s.Id = id
}

func SelfAddress(s SelfAddressable) error {
	s.SetId("############################################")

	message, err := json.Marshal(s)
	if err != nil {
		return err
	}

	buffer := [33]byte{}
	sum := blake3.Sum256(message)
	copy(buffer[1:], sum[:])

	b64 := base64.URLEncoding.EncodeToString(buffer[:])
	qb64 := []rune(b64)
	qb64[0] = 'E'

	s.SetId(string(qb64))

	return nil
}

func VerifyAddress(s SelfAddressable) error {
	oldId := s.GetId()

	if err := SelfAddress(s); err != nil {
		return err
	}

	if !strings.EqualFold(s.GetId(), oldId) {
		return fmt.Errorf("address verification failed")
	}

	s.SetId(oldId)

	return nil
}
