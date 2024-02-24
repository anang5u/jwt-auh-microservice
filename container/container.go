package container

import (
	"jwt-auth-microservice/config"
	"jwt-auth-microservice/infrastructure/auth"
	"jwt-auth-microservice/infrastructure/logger"
	"jwt-auth-microservice/uses-cases/authuc"
	"log"
)

type Container struct {
	EnvironmentConfig config.EnvironmentConfig
	Auth              *authuc.JwtAuthUseCase
}

func SetupContainer() Container {
	log.Println("Starting new container...")

	// Setup config
	log.Println("Setup config...")
	config, err := config.SetupConfig()
	if err != nil {
		log.Panic(err)
	}

	// Setup logger
	log.Println("Setup logger...")
	logger.SetupLogger(config)

	// Loading infrastructures
	log.Println("Loading infrastructure's...")
	authManager := auth.NewJwtAuth()

	// Loading uses case (service)
	log.Println("Loading uses case's...")
	authUseCase := authuc.NewJwtAuthUseCase(authManager)

	// Init http handler's
	log.Println("Init http handler's...")

	return Container{
		EnvironmentConfig: config,
		Auth:              authUseCase,
	}
}
