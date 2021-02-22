package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"cleanrss/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type ErrorResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

//go:embed static/*
var static embed.FS

var Server *fiber.App

func ServerInit(){
	Server = fiber.New(fiber.Config{DisableStartupMessage: true})
	Server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	ServeStatic()
	apiGroup := Server.Group("/api")
	routes.RoutesInit(apiGroup)

	log.Println("Server will listen on http://localhost:1337")
	Server.Listen("localhost:1337")
}

func ServeStatic() {
	staticRouteFixed, _ := fs.Sub(static, "static")
	Server.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(staticRouteFixed),
	}))
}

func ErrorResponseFactory(httpCode int, errCode string, err error, c *fiber.Ctx) error {
	c.Status(400)
	return c.JSON(ErrorResponse{errCode, err.Error()})
}