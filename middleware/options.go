package middleware

import (
	"github.com/52dev/kratos-authn/engine"
	"github.com/go-kratos/kratos/v2/log"
)

type Option func(*options)

type options struct {
	claims engine.AuthClaims
	log    *log.Helper
}

func WithAuthClaims(claims engine.AuthClaims) Option {
	return func(o *options) {
		o.claims = claims
	}
}

func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.log = log.NewHelper(log.With(logger, "module", "authn.middleware"))
	}
}
