package outbound

type SessionTokenHasher interface {
	Hash(token string) (string, error)
}