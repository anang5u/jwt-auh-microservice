package authuc

import "jwt-auth-microservice/domain/entities"

type authManager interface {
	GenerateRSAPrivKey() error
	GetJwks() (*entities.JSONWebKeySet, error)
	GenerateToken() (*entities.LoginResponse, error)
}
