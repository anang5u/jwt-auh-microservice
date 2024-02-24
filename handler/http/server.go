package httphandler

import (
	cnt "jwt-auth-microservice/container"

	"github.com/gin-gonic/gin"
)

func ServeHttp(container cnt.Container) {
	router := gin.Default()

	setupRouteV1(router, container)
	router.Run(":" + container.EnvironmentConfig.App.Port)
}
