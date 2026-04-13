package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/52dev/kratos-authn/engine"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

func (hc headerCarrier) Set(key, value string) {
	http.Header(hc).Set(key, value)
}

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Add append value to key-values pair.
func (hc headerCarrier) Add(key string, value string) {
	http.Header(hc).Add(key, value)
}

// Values returns a slice of values associated with the passed key.
func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}

type myTransporter struct {
	reqHeader   headerCarrier
	replyHeader headerCarrier
}

func (t *myTransporter) Kind() transport.Kind            { return "test" }
func (t *myTransporter) Endpoint() string                { return "" }
func (t *myTransporter) Operation() string               { return "" }
func (t *myTransporter) RequestHeader() transport.Header { return t.reqHeader }
func (t *myTransporter) ReplyHeader() transport.Header   { return t.replyHeader }

func TestAuthenticator(t *testing.T) {
	ctx := context.Background()

	ctx = transport.NewServerContext(ctx, &myTransporter{reqHeader: headerCarrier{}, replyHeader: headerCarrier{}})
	ctx = transport.NewClientContext(ctx, &myTransporter{reqHeader: headerCarrier{}, replyHeader: headerCarrier{}})

	auth, err := NewAuthenticator(
		WithKey([]byte("test")),
		WithSigningMethod("HS256"),
	)
	assert.Nil(t, err)

	scopes := []string{"local:admin:user_name", "tenant:admin:user_name"}

	principal := engine.AuthClaims{
		engine.ClaimFieldSubject: "user_name",
		engine.ClaimFieldScope:   scopes,
	}

	outToken, err := auth.CreateIdentity(principal)
	assert.Nil(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6WyJsb2NhbDphZG1pbjp1c2VyX25hbWUiLCJ0ZW5hbnQ6YWRtaW46dXNlcl9uYW1lIl0sInN1YiI6InVzZXJfbmFtZSJ9.xIzbQbQSlzdms5ZVaHrg6pZohDlt0DTYopobUo2qqQw", outToken)

	ctx, err = auth.CreateIdentityWithContext(ctx, engine.ContextTypeKratosMetaData, principal)
	assert.Nil(t, err)

	var token string
	if header, ok := transport.FromClientContext(ctx); ok {
		str := header.RequestHeader().Get("Authorization")
		splits := strings.SplitN(str, " ", 2)
		assert.Equal(t, 2, len(splits))
		assert.Equal(t, engine.BearerWord, splits[0])
		token = str
		//fmt.Println(token)
	}

	if header, ok := transport.FromServerContext(ctx); ok {
		header.RequestHeader().Set("Authorization", token)
	}

	authToken, err := auth.Authenticate(ctx, engine.ContextTypeKratosMetaData)
	assert.Nil(t, err)

	sub, _ := authToken.GetSubject()
	assert.Equal(t, "user_name", sub)

	scopesOut, _ := authToken.GetScopes()
	assert.Equal(t, 2, len(scopesOut))
	assert.Equal(t, "local:admin:user_name", scopesOut[0])
	assert.Equal(t, "tenant:admin:user_name", scopesOut[1])
	fmt.Println(authToken)
}

func TestRS256Authenticator(t *testing.T) {
	// 加载 RSA 私钥
	privKey, err := jwtV5.ParseRSAPrivateKeyFromPEM([]byte("-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDXLNKj9V7iOeUc\nRbx4AswYnAkPjn57F/Z9OaJHdRCPmca53bOzrOQCSIAOkvtsQPOA1oPAfttsnL4G\nHlvQIep+/TnNq0IYpmkzghGLZTiyK9scMdqVBc636vU8pB+8nyrYCeqTPJEBUpOH\n0IIqWg7EaieOrhNbGm+2Yqe2dg5R7NmwHaI4+Pky8ISLl2YjlNgWW3r/tbFCjJBg\nCB9DHfeY7jBfp2MrMzXztgKU+ff0fnoQjZdma+RAodFwunaY0B3s1Zc2ta+eu3eb\nYeZ3ibrTMG6TdWy8QFVhM75rdei6UOdYUS+pYYsAG9/4a2flrKUzhjV3zDSqpmtK\nZLUEneGdAgMBAAECggEAH2O+DMYYQ+dPOXsg3d8GmBZ3KepPIDTkM+tq9YKp2lEE\nEQw7EVyI3J5n8/hULjwhaauhh7zZ1LPe8rSOD0RLWaAmRQ8VMtRf53AzkALBrRhB\nvBC3wuKYf/MKOID20kTj8qUrr7P3sVozBG6R9oyxt8yGncVeNH2cS16D+dWqDCA9\nMRjL1PLgv5HxcdhjL4m8EMt1zmRzXMpDqflpjlavbj7HXqmZRFFj0qf+aYkhse51\nUpum/QEu9nyYsgQS+CgqAZ5kZrZBUaUsWp3TFBrRGq66FvFrEn0QfFMEORn+GjNB\nsb+qObvIKWVaNp2OywDpKWe4R45Tk9/Oi2ov9gmWbwKBgQDzneJeRIZFFMXJObno\nhK33aL+SGTdr58ONHQJEwce0cT64d7uLuGcoHhukBhnBMB/b++I8+UOBnHzmMZp2\nEXeK5OKBEyErOIEaX/ltsUl3AycRMQ5N0GA/HZPGcJATnVlB1ja5rOocCrjZbPiw\nhAR3XvgfcDsB7g7RkAOX0fX3iwKBgQDiHNXgeOu+Ns9UHTaJSLWXk5unFjq/DIbI\nDUQtAudVAO8IlTHujx9u8Hp9XTQgWmjW5kCa0gxdPbig5c27CV7qhE3OTqxcc3FM\nAb1IPkfxum6kEhlOj1S21aAkNlYrVFUPc13Rl62lXURnQAmKEssvcRvnl1RG/O04\nvIPp9FFwdwKBgCHvRz4MW4u55gcutFfQS49gFvdZ7d9pDFNWzB8ZwyC+eZcmjoha\n6nurHfyOIP5JHtb80jneGuouCzPhivuRWU6OrYJ/UKp9l3Y+EjeWb35VgRai97Qd\nJ5sDGreUrG0fCPTjywG4NXAsii03QbkM2rZqEzQF5SJSr9u/LND0HUgbAoGBAJWp\nJF9JajAisyQnmdtQRvGm/9WePxAJSITNUxy/2UJINe7mYYBXNyUFAu5LbJ8leFMV\nYBmZghmNKtFEieGMmEh9fcpaBHfE6W63kANrRc9X6LesSxfWgunph++wD2TqksqB\nP83kqUjU7NuyZR4AxoAGS8QERAIgkxuEm4OU9PqNAoGBANXnfGy87HXD7iILc43O\nJNDJA7gPY3+KwNueL+fO2jXRnIsrbv6hTNuatJ8twtV4AJTPFXgARa3zlRCfHaCo\npPqXP2c1Fyb26mD9HVBA3PRLMs3IOnhgzfptYJqUANJlC+6VUVVW7gVcoJr7d4ik\nbWCwfqfzRh3DLFGEwBiuaJYH\n-----END PRIVATE KEY-----\n"))

	// 加载 RSA 公钥
	pubKey, err := jwtV5.ParseRSAPublicKeyFromPEM([]byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1yzSo/Ve4jnlHEW8eALM\nGJwJD45+exf2fTmiR3UQj5nGud2zs6zkAkiADpL7bEDzgNaDwH7bbJy+Bh5b0CHq\nfv05zatCGKZpM4IRi2U4sivbHDHalQXOt+r1PKQfvJ8q2AnqkzyRAVKTh9CCKloO\nxGonjq4TWxpvtmKntnYOUezZsB2iOPj5MvCEi5dmI5TYFlt6/7WxQoyQYAgfQx33\nmO4wX6djKzM187YClPn39H56EI2XZmvkQKHRcLp2mNAd7NWXNrWvnrt3m2Hmd4m6\n0zBuk3VsvEBVYTO+a3XoulDnWFEvqWGLABvf+Gtn5aylM4Y1d8w0qqZrSmS1BJ3h\nnQIDAQAB\n-----END PUBLIC KEY-----\n"))
	ctx := context.Background()

	ctx = transport.NewServerContext(ctx, &myTransporter{reqHeader: headerCarrier{}, replyHeader: headerCarrier{}})
	ctx = transport.NewClientContext(ctx, &myTransporter{reqHeader: headerCarrier{}, replyHeader: headerCarrier{}})

	auth, err := NewAuthenticator(
		WithKey([]byte("test")),
		WithSigningMethod("RS256"),
		WithRSAKeys(privKey, pubKey),
	)
	assert.Nil(t, err)

	scopes := []string{"local:admin:user_name", "tenant:admin:user_name"}

	principal := engine.AuthClaims{
		engine.ClaimFieldSubject: "user_name",
		engine.ClaimFieldScope:   scopes,
	}

	outToken, err := auth.CreateIdentity(principal)
	assert.Nil(t, err)
	assert.Equal(t, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6WyJsb2NhbDphZG1pbjp1c2VyX25hbWUiLCJ0ZW5hbnQ6YWRtaW46dXNlcl9uYW1lIl0sInN1YiI6InVzZXJfbmFtZSJ9.LATfIBcDfOuCo1AOWUvk_UwPepO-2g6TAsOUE_KZYDGaJJVm7R9RPqTRgEnOrTg1obT7uCzm13quSnMhYBVkDHeb10Vin4A1GAxcLrh_SeK7c_A5EE_7MNN6h-Ib5fAxwlE_K0GVVxXVfRKsjMhyeHtUhW0GG5xOJXqUkjNBPORCpY8IXUmhZ_4L-GN-pt4Hs7HQXn25Dap-0HVUnhfjDJ_X-cbehVnR4iDNoaVx3_E3f_Ab2UreWNzs1yDTeaW8d4TEq8S6C8RgpLoY7kxOsuaa2B8d1HnRDz-le6CP98_CqS7XC9S61vx3f0Cqf0fv-Aq8UYerUNg0Go9sBCvAbw", outToken)

	c, err := auth.AuthenticateToken(outToken)
	assert.Nil(t, err)
	subject, err := c.GetSubject()
	assert.Nil(t, err)
	assert.Equal(t, "user_name", subject)
	ctx, err = auth.CreateIdentityWithContext(ctx, engine.ContextTypeKratosMetaData, principal)
	assert.Nil(t, err)

	var token string
	if header, ok := transport.FromClientContext(ctx); ok {
		str := header.RequestHeader().Get("Authorization")
		splits := strings.SplitN(str, " ", 2)
		assert.Equal(t, 2, len(splits))
		assert.Equal(t, engine.BearerWord, splits[0])
		token = str
		//fmt.Println(token)
	}

	if header, ok := transport.FromServerContext(ctx); ok {
		header.RequestHeader().Set("Authorization", token)
	}

	authToken, err := auth.Authenticate(ctx, engine.ContextTypeKratosMetaData)
	assert.Nil(t, err)

	sub, _ := authToken.GetSubject()
	assert.Equal(t, "user_name", sub)

	scopesOut, _ := authToken.GetScopes()
	assert.Equal(t, 2, len(scopesOut))
	assert.Equal(t, "local:admin:user_name", scopesOut[0])
	assert.Equal(t, "tenant:admin:user_name", scopesOut[1])
	fmt.Println(authToken)
}
