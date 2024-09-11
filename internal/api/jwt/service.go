package jwt

import (
	"crypto/ed25519"
	"time"

	"github.com/o1egl/paseto"
)

type JWTService struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func InitKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func (j *JWTService) GetToken(key, value string, expirationTime time.Duration) (*string, error) {
	jsonToken := paseto.JSONToken{
		Expiration: time.Now().Add(expirationTime),
	}
	jsonToken.Set(key, value)
	token, err := paseto.NewV2().Sign(j.privateKey, jsonToken, nil)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (j *JWTService) Verify(token string) (*paseto.JSONToken, error) {
	var parsedToken paseto.JSONToken
	err := paseto.NewV2().Verify(token, j.publicKey, &parsedToken, nil)
	if err != nil {
		return nil, err
	}
	return &parsedToken, nil
}

func InitJWTService(publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) *JWTService {
	return &JWTService{publicKey: publicKey, privateKey: privateKey}
}
