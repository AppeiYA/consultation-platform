package uuid

import (
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
	"github.com/google/uuid"
)

var errPrefixIsRequired = custom_errors.InternalServerError("no prefix sent to generator")

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(prefix string) (string, error) {
	if prefix == "" {
		return "", errPrefixIsRequired
	}
	return prefix + "_" + uuid.NewString(), nil
}