package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	"github.com/AFK068/skilloom/ingest/internal/services/auth"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/respond"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/authenticator_mock.go -package=mock

const principalKey = "principal"

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (principal.Principal, error)
}

type AuthDeps struct {
	fx.In

	Auth Authenticator
	Log  *zap.Logger
}

type Auth struct {
	AuthDeps
	log *zap.Logger
}

func NewAuth(d AuthDeps) *Auth {
	return &Auth{
		AuthDeps: d,
		log:      d.Log.Named("auth middleware"),
	}
}

func (a *Auth) Handle(c *gin.Context) {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok || token == "" {
		respond.Error(c, http.StatusUnauthorized, "missing or invalid bearer token")
		return
	}

	p, err := a.Auth.Authenticate(c.Request.Context(), token)
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		respond.Error(c, http.StatusUnauthorized, "invalid token")
		return
	case err != nil:
		a.log.Error("authenticate failed", zap.Error(err))
		respond.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.Set(principalKey, p)
	c.Next()
}

func Principal(c *gin.Context) (principal.Principal, bool) {
	val, _ := c.Get(principalKey)
	p, ok := val.(principal.Principal)
	return p, ok
}
