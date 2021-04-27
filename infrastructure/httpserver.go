package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"sync"
)

type HTTPServer struct {
	*fiber.App
}

func NewHTTPServer() HTTPServer {
	var server HTTPServer
	server.App = fiber.New(fiber.Config{DisableStartupMessage: true})
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	return server
}

func (server HTTPServer) Listen(addr string, wg *sync.WaitGroup) error {
	defer wg.Done()
	return server.App.Listen(addr)
}
