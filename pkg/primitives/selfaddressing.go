package primitives

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
