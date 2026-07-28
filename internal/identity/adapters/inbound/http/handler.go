package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
)

type IdentityHandler struct {
	identityModule *identity.Module
	cookieManager shared_http.CookieManagerInt
}

func NewIdentityHandler(module *identity.Module, cookieManager shared_http.CookieManagerInt) *IdentityHandler {
	return &IdentityHandler{
		identityModule: module,
		cookieManager: cookieManager,
	}
}

