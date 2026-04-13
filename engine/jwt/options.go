package jwt

import (
	"crypto/rsa"

	jwtV5 "github.com/golang-jwt/jwt/v5"
)

type Options struct {
	signingMethod jwtV5.SigningMethod
	keyFunc       jwtV5.Keyfunc

	// 新增：RSA 专用
	privateKey *rsa.PrivateKey // 签名用
	publicKey  *rsa.PublicKey  // 验证用
}

type Option func(d *Options)

// WithSigningMethod set signing method
func WithSigningMethod(alg string) Option {
	return func(o *Options) {
		o.signingMethod = jwtV5.GetSigningMethod(alg)
	}
}

// WithKey set key
func WithKey(key interface{}) Option {
	return func(o *Options) {
		o.keyFunc = func(token *jwtV5.Token) (interface{}, error) {
			return key, nil
		}
	}
}

// WithRSAKeys 设置 RSA 公私钥（推荐）
func WithRSAKeys(priv *rsa.PrivateKey, pub *rsa.PublicKey) Option {
	return func(o *Options) {
		o.privateKey = priv
		o.publicKey = pub

		// 自动设置 keyFunc：验证时返回公钥
		o.keyFunc = func(token *jwtV5.Token) (interface{}, error) {
			return pub, nil
		}
	}
}

// WithRSAPrivateKey 只设置私钥（适合只签发不验证）
func WithRSAPrivateKey(priv *rsa.PrivateKey) Option {
	return func(o *Options) {
		o.privateKey = priv
		o.keyFunc = func(token *jwtV5.Token) (interface{}, error) {
			return &priv.PublicKey, nil
		}
	}
}
