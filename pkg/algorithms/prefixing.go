package algorithms

import (
	"fmt"
	"strings"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

func CreatePrefix(p primitives.Prefixable) error {
	p.SetPrefix("############################################")

	if err := SelfAddress(p); err != nil {
		return err
	}

	p.SetPrefix(p.GetId())

	return nil
}

func VerifyPrefixAndData(p primitives.Prefixable) error {
	oldId := p.GetId()
	oldPrefix := p.GetPrefix()

	defer func() {
		p.SetId(oldId)
		p.SetPrefix(oldPrefix)
	}()

	if err := CreatePrefix(p); err != nil {
		return err
	}

	if !strings.EqualFold(p.GetId(), oldId) {
		return fmt.Errorf("address verification failed")
	}

	if !strings.EqualFold(p.GetId(), oldPrefix) {
		return fmt.Errorf("prefix verification failed")
	}

	return nil

}
