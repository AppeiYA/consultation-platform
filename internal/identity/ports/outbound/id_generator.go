package outbound

type IdentifierGenerator interface {
	Generate() (string, error)
}