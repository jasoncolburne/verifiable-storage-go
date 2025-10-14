package verifiablestorage

import (
	"fmt"
	"strings"
)

type Prefixable interface {
	SelfAddressable
	GetPrefix() string
	SetPrefix(prefix string)
}

type Prefixer struct {
	SelfAddresser
	Prefix string `db:"prefix" json:"prefix"`
}

func (p Prefixer) GetPrefix() string {
	return p.Prefix
}

func (p *Prefixer) SetPrefix(prefix string) {
	p.Prefix = prefix
}

func CreatePrefix(p Prefixable) error {
	p.SetPrefix("############################################")

	if err := SelfAddress(p); err != nil {
		return err
	}

	p.SetPrefix(p.GetId())

	return nil
}

func VerifyPrefix(p Prefixable) error {
	oldId := p.GetId()
	oldPrefix := p.GetPrefix()

	if err := CreatePrefix(p); err != nil {
		return err
	}

	if !strings.EqualFold(p.GetId(), oldId) {
		return fmt.Errorf("address verification failed")
	}

	if !strings.EqualFold(p.GetId(), oldPrefix) {
		return fmt.Errorf("prefix verification failed")
	}

	p.SetId(oldId)
	p.SetPrefix(oldPrefix)

	return nil

}
