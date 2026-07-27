package http

import "github.com/AppeiYA/consultation-platform/internal/identity"

type IdentityHandler struct {
	identityModule *identity.Module
}

func NewIdentityHandler(module *identity.Module) *IdentityHandler {
	return &IdentityHandler{
		identityModule: module,
	}
}

