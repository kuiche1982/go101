package income

type Income interface {
	Calculate() int
	Source() string
}
