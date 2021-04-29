package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"sync"
)

type HTTPFiberServer struct {
	*fiber.App
}

func NewHTTPFiberServer() HTTPFiberServer {
	var server HTTPFiberServer
	server.App = fiber.New(fiber.Config{DisableStartupMessage: true})
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	return server
}

func (server HTTPFiberServer) Listen(addr string, wg *sync.WaitGroup) error {
	defer wg.Done()
	return server.App.Listen(addr)
}
