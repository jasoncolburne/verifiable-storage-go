package crypto

type Noncer interface {
	Generate() (string, error)
}
