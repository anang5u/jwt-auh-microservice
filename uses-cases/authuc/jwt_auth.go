package authuc

import "jwt-auth-microservice/domain/entities"

type JwtAuthUseCase struct {
	authManager authManager
}

func NewJwtAuthUseCase(am authManager) *JwtAuthUseCase {
	return &JwtAuthUseCase{
		authManager: am,
	}
}

func (uc *JwtAuthUseCase) Login() (*entities.LoginResponse, error) {
	return uc.authManager.GenerateToken()
}

func (uc *JwtAuthUseCase) GetJwks() (*entities.JSONWebKeySet, error) {
	return uc.authManager.GetJwks()
}
