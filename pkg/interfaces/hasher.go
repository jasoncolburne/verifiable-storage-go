package interfaces

type Hasher interface {
	Sum(message string) string
}
