package main

import (
	"AuthService/pkg/handler"
	"AuthService/pkg/server"
	"AuthService/pkg/service"
	"AuthService/pkg/storage"
)

// @title Auth Service API
// main godoc
func main() {
	repository := storage.NewPostgreSQLStorage()
	services := service.NewService(repository)
	api := handler.NewHandler(services)

	srv := server.NewServer()
	srv.Run(api.InitRoutes())

}
