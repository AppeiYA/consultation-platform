package http

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"github.com/gofiber/fiber/v2"
)

type CookieManagerInt interface {
	SetSession(c *fiber.Ctx, token string)
	GetSession(c *fiber.Ctx) (string, bool)
	DeleteSession(c *fiber.Ctx)
}

type CookieManager struct {
	cfg config.SessionConfig
}

func NewCookieManager(cfg config.SessionConfig) CookieManagerInt {
	return &CookieManager{cfg: cfg}
}

func (cm *CookieManager) SetSession(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     cm.cfg.CookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: cm.cfg.HTTPOnly,
		Secure:   cm.cfg.Secure,
		SameSite: cm.cfg.SameSite,
		MaxAge:   int(cm.cfg.TTL.Seconds()),
		Expires: time.Now().Add(cm.cfg.TTL),
	})
}

func (cm *CookieManager) GetSession(c *fiber.Ctx) (string, bool) {
    token := c.Cookies(cm.cfg.CookieName)
    if token == "" {
        return "", false
    }

    return token, true
}

func (cm *CookieManager) DeleteSession(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     cm.cfg.CookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: cm.cfg.HTTPOnly,
		Secure:   cm.cfg.Secure,
		SameSite: cm.cfg.SameSite,
		MaxAge:   -1,
		Expires: time.Unix(0,0),
	})
}

