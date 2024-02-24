package httphandler

import (
	cnt "jwt-auth-microservice/container"
	authhandler "jwt-auth-microservice/handler/http/auth"
)

type httpHandler struct {
	auth *authhandler.AuthHandler
}

func NewHttpHandler(container cnt.Container) *httpHandler {
	return &httpHandler{
		auth: authhandler.NewAuthHandler(container),
	}
}
