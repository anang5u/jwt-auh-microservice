package main

import (
	"jwt-auth-microservice/container"
	httphandler "jwt-auth-microservice/handler/http"
)

func main() {
	// start init container
	container := container.SetupContainer()

	// serv http handler
	httphandler.ServeHttp(container)
}
