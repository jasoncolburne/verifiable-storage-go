package primitives

// This file is a bit different. Constraints are composed by the model, and this is simply
// for convenience

type Searchable interface {
	GetSearchKey() string
	SetSearchKey(hash string)
}

type Searcher struct {
	SearchKey string `db:"search_key" json:"searchKey"`
}

func (s Searcher) GetSearchKey() string {
	return s.SearchKey
}

func (s *Searcher) SetSearchKey(key string) {
	s.SearchKey = key
}
