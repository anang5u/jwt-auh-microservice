package authhandler

import (
	cnt "jwt-auth-microservice/container"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	container cnt.Container
}

func NewAuthHandler(c cnt.Container) *AuthHandler {
	return &AuthHandler{
		container: c,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	response, err := h.container.Auth.Login()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed. Err: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetJwks(c *gin.Context) {
	response, err := h.container.Auth.GetJwks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get jwks failed. Err: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
