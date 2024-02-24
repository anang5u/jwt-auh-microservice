package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"jwt-auth-microservice/domain/entities"
	"jwt-auth-microservice/shared"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type jwtAuth struct {
	privKeyFilename  string
	jwksJsonFilename string
}

// NewJwtAuth
func NewJwtAuth() *jwtAuth {
	return &jwtAuth{
		privKeyFilename:  JWTPrivKeyFilename,
		jwksJsonFilename: JWKSJsonFilename,
	}
}

// SetPrivateKeyFilename
func (m *jwtAuth) SetPrivateKeyFilename(filename string) *jwtAuth {
	m.privKeyFilename = filename
	return m
}

// SetJwksJsonFilename
func (m *jwtAuth) SetJwksJsonFilename(filename string) *jwtAuth {
	m.jwksJsonFilename = filename
	return m
}

// GenerateRSAPrivKey
func (m *jwtAuth) GenerateRSAPrivKey() error {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Marshal the private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFile, err := os.Create(m.privKeyFilename)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	// Write the private key to the file
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	return nil
}

// GetJwks
func (m *jwtAuth) GetJwks() (*entities.JSONWebKeySet, error) {
	if m.isJwksJsonExists() {
		jwkSet := entities.JSONWebKeySet{}
		err := shared.JSONFromFile(m.jwksJsonFilename, &jwkSet)
		if err != nil {
			return nil, err
		}

		return &jwkSet, nil
	}

	// generate
	jwks, err := m.generateJwks()
	if err != nil {
		return nil, err
	}

	// write to json file
	err = shared.FileWriteToFile(jwks, m.jwksJsonFilename)
	if err != nil {
		return nil, err
	}

	return jwks, nil
}

// GenerateToken
func (m *jwtAuth) GenerateToken() (*entities.LoginResponse, error) {
	// JWK
	jwks, err := m.GetJwks()
	if err != nil {
		return nil, err
	}
	keys := jwks.Keys[0]

	// jwkByte
	jwkByte, err := json.Marshal(&keys)
	if err != nil {
		return nil, err
	}
	jwkStr := string(jwkByte)

	// Build a JWT!
	token := jwt.New()
	token.Set(jwt.SubjectKey, JWTSubjectKey)
	token.Set(jwt.IssuerKey, JWTIssuerKey)
	token.Set(jwt.IssuedAtKey, time.Unix(500, 0))

	jwkey, err := jwk.ParseKey([]byte(jwkStr))
	if err != nil {
		return nil, err
	}

	// signed and return a jwt
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, jwkey))
	if err != nil {
		return nil, err
	}

	result := entities.LoginResponse{
		Token: string(signed),
		User: entities.UserResponse{
			ID: "dummy-user-id",
		},
	}

	return &result, nil
}

// isJwksJsonExists
func (m *jwtAuth) isJwksJsonExists() bool {
	if _, err := os.Stat(m.jwksJsonFilename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// generateJwks
func (m *jwtAuth) generateJwks() (*entities.JSONWebKeySet, error) {
	key, err := jwk.ReadFile(m.privKeyFilename, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	// Mengkonversi *jose.JSONWebKey ke JSON
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}

	// Deklarasikan objek jwk.Set
	jwkSet := entities.JSONWebKeySet{}

	// Lakukan Unmarshal JSON ke dalam objek jwk.Set
	if err := json.Unmarshal(keyJSON, &jwkSet); err != nil {
		return nil, err
	}

	return &jwkSet, nil
}
