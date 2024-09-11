package jwt

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	t.Run("InitKey() generate public,private key properly", func(t *testing.T) {
		publicKey, privateKey, err := InitKey()
		if err != nil {
			t.Fatalf("Fail to create key pair: %v", err)
		}
		assert.Equal(t, len(publicKey), ed25519.PublicKeySize)
		assert.Equal(t, len(privateKey), ed25519.PrivateKeySize)
	})

	t.Run("GetAccessKey generate token properly", func(t *testing.T) {
		publicKey, privateKey, err := InitKey()
		if err != nil {
			t.Fatalf("Fail to create key pair: %v", err)
		}
		jwtService := InitJWTService(publicKey, privateKey)
		key := "id"
		value := "1"
		accessToken, err := jwtService.GetToken(key, value, ACCESSTOKEN_EXPIRATION_TIME)
		if err != nil {
			t.Fatalf("Fail to create token: %v", err)
		}

		refreshToken, err := jwtService.GetToken(key, value, ACCESSTOKEN_EXPIRATION_TIME)
		if err != nil {
			t.Fatalf("Fail to create token: %v", err)
		}

		var parsedAccessToken paseto.JSONToken
		err = paseto.NewV2().Verify(*accessToken, publicKey, &parsedAccessToken, nil)
		if err != nil {
			t.Fatalf("Fail to verify token: %v", err)
		}
		id := parsedAccessToken.Get(key)
		expirationTime := parsedAccessToken.Expiration

		assert.Equal(t, value, id)
		assert.False(t, expirationTime.Before(time.Now()))

		var parsedRefreshToken paseto.JSONToken
		err = paseto.NewV2().Verify(*refreshToken, publicKey, &parsedRefreshToken, nil)
		if err != nil {
			t.Fatalf("Fail to verify token: %v", err)
		}
		id = parsedRefreshToken.Get(key)
		expirationTime = parsedRefreshToken.Expiration

		assert.Equal(t, value, id)
		assert.False(t, expirationTime.Before(time.Now()))
	})

	t.Run("Verify return token when valid token is given", func(t *testing.T) {
		publicKey, privateKey, err := InitKey()
		if err != nil {
			t.Fatalf("Fail to create key pair: %v", err)
		}
		jwtService := InitJWTService(publicKey, privateKey)
		token, err := jwtService.GetToken("some", "value", ACCESSTOKEN_EXPIRATION_TIME)
		if err != nil {
			t.Fatalf("Token generation fail: %v", err)
		}
		parsedToken, err := jwtService.Verify(*token)
		if err != nil {
			t.Fatalf("Token verify fail: %v", err)
		}
		value := parsedToken.Get("some")
		assert.Equal(t, "value", value)
	})

	t.Run("Verify return error when invalid token is given", func(t *testing.T) {
		publicKey, privateKey, err := InitKey()
		if err != nil {
			t.Fatalf("Fail to create key pair: %v", err)
		}
		jwtService := InitJWTService(publicKey, privateKey)

		_, err = jwtService.Verify("some invalid token")
		assert.Error(t, err)
	})
}
