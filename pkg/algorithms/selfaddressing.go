package algorithms

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
	"github.com/zeebo/blake3"
)

func SelfAddress(s primitives.SelfAddressable) error {
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

func VerifyAddress(s primitives.SelfAddressable) error {
	oldId := s.GetId()

	defer func() {
		s.SetId(oldId)
	}()

	if err := SelfAddress(s); err != nil {
		return err
	}

	if !strings.EqualFold(s.GetId(), oldId) {
		return fmt.Errorf("address verification failed")
	}

	return nil
}
