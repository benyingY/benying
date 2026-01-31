package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"enterprise-llm-agent-platform/shared"
	"github.com/gin-gonic/gin"
)

type authClaims struct {
	Subject   string `json:"sub"`
	Issuer    string `json:"iss,omitempty"`
	Audience  any    `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

var publicKeyCache struct {
	mu   sync.RWMutex
	pem  string
	key  *rsa.PublicKey
	err  error
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := shared.GetRequestID(c)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, errorResponse{
				Status:    "error",
				Service:   serviceName,
				RequestID: requestID,
				Error:     "unauthorized",
				Message:   "missing bearer token",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, errorResponse{
				Status:    "error",
				Service:   serviceName,
				RequestID: requestID,
				Error:     "unauthorized",
				Message:   "missing bearer token",
			})
			c.Abort()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse{
				Status:    "error",
				Service:   serviceName,
				RequestID: requestID,
				Error:     "unauthorized",
				Message:   "invalid token",
			})
			c.Abort()
			return
		}

		tenantID, _ := resolveTenantID(*claims)
		c.Set("tenant_id", tenantID)
		c.Set("user_id", claims.Subject)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func parseToken(tokenString string) (*authClaims, error) {
	key, err := getPublicKey()
	if err != nil {
		return nil, err
	}
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid header encoding")
	}
	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errors.New("invalid header")
	}
	if header.Alg != "RS256" {
		return nil, errors.New("unsupported alg")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}
	var claims authClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims")
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid signature encoding")
	}
	signingInput := parts[0] + "." + parts[1]
	hashed := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature); err != nil {
		return nil, errors.New("invalid signature")
	}

	if claims.ExpiresAt != 0 && time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token expired")
	}
	if issuer := os.Getenv("ACCESS_JWT_ISSUER"); issuer != "" && claims.Issuer != issuer {
		return nil, errors.New("invalid issuer")
	}
	if audience := os.Getenv("ACCESS_JWT_AUDIENCE"); audience != "" && !audienceMatches(claims.Audience, audience) {
		return nil, errors.New("invalid audience")
	}

	return &claims, nil
}

func getPublicKey() (*rsa.PublicKey, error) {
	pem := os.Getenv("ACCESS_JWT_PUBLIC_KEY")
	if pem == "" {
		return nil, errors.New("missing public key")
	}

	publicKeyCache.mu.RLock()
	if publicKeyCache.pem == pem && publicKeyCache.key != nil {
		key := publicKeyCache.key
		publicKeyCache.mu.RUnlock()
		return key, nil
	}
	if publicKeyCache.pem == pem && publicKeyCache.err != nil {
		err := publicKeyCache.err
		publicKeyCache.mu.RUnlock()
		return nil, err
	}
	publicKeyCache.mu.RUnlock()

	key, err := parsePublicKeyPEM([]byte(pem))
	publicKeyCache.mu.Lock()
	publicKeyCache.pem = pem
	publicKeyCache.key = key
	publicKeyCache.err = err
	publicKeyCache.mu.Unlock()
	return key, err
}

func parsePublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid public key")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		if key, ok := parsedKey.(*rsa.PublicKey); ok {
			return key, nil
		}
		return nil, errors.New("unsupported public key type")
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("invalid public key")
	}
	return key, nil
}

func audienceMatches(audience any, expected string) bool {
	switch value := audience.(type) {
	case string:
		return value == expected
	case []any:
		for _, item := range value {
			if str, ok := item.(string); ok && str == expected {
				return true
			}
		}
	}
	return false
}

// resolveTenantID is intentionally minimal for now; fill in mapping rules later.
func resolveTenantID(claims authClaims) (string, error) {
	return "", nil
}
