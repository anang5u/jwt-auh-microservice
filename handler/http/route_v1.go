package httphandler

import (
	cnt "jwt-auth-microservice/container"

	"github.com/gin-gonic/gin"
)

func setupRouteV1(router *gin.Engine, container cnt.Container) {
	httpHandler := NewHttpHandler(container)

	v1 := router.Group("/v1")
	{
		v1.POST("/login", httpHandler.auth.Login)
		v1.GET("/jwks", httpHandler.auth.GetJwks)
	}

}
