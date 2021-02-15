package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

var Server *fiber.App

func ServerInit(){
	Server = fiber.New(fiber.Config{DisableStartupMessage: true})
	RoutesInit(Server)
	log.Println("Server will listen on http://localhost:1337")
	Server.Listen("localhost:1337")
}

func ErrorResponseFactory(httpCode int, errCode string, err error, c *fiber.Ctx) error {
	c.Status(400)
	return c.JSON(ErrorResponse{errCode, err.Error()})
}