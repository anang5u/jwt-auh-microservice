package auth_test

import (
	"jwt-auth-microservice/infrastructure/auth"
	"jwt-auth-microservice/shared"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	privKeyFile  = "../../certs/private_key.pem"
	jwksJsonFile = "../../jwks_test.json"
)

func TestJwtAuth_GenerateRSAPrivKey(t *testing.T) {
	jwtAuth := auth.NewJwtAuth()
	err := jwtAuth.SetPrivateKeyFilename("../../certs/private_key_test.pem").GenerateRSAPrivKey()

	assert.Nil(t, err, "error should be nil")
}

func TestJwtAuth_GetJwks(t *testing.T) {
	jwtAuth := auth.NewJwtAuth()
	jwks, err := jwtAuth.SetPrivateKeyFilename(privKeyFile).SetJwksJsonFilename(jwksJsonFile).GetJwks()

	log.Println("jwks >>>")
	shared.JSONPretty(jwks)

	assert.Nil(t, err, "error should be nil")
}

func TestJwtAuth_GenerateToken(t *testing.T) {
	jwtAuth := auth.NewJwtAuth()
	tokenStr, err := jwtAuth.SetPrivateKeyFilename(privKeyFile).SetJwksJsonFilename(jwksJsonFile).GenerateToken()
	if tokenStr != nil {
		log.Println("token >>>")
		log.Println(*tokenStr)
	}

	assert.Nil(t, err, "error should be nil")
}
