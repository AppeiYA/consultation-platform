package outbound

type IdentifierGenerator interface {
	Generate(prefix string) (string, error)
}