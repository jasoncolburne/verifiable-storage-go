package interfaces

type Noncer interface {
	Generate() (string, error)
}
