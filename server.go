package main

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

var Server *fiber.App

func ServerInit(){
	Server = fiber.New()
	RoutesInit(Server)
	Server.Listen(":1337")
}

func ErrorResponseFactory(httpCode int, errCode string, err error, c *fiber.Ctx) error {
	c.Status(400)
	return c.JSON(ErrorResponse{errCode, err.Error()})
}