package routes

import (
	"cleanrss/controllers/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func wsRouter(router fiber.Router) {
	go ws.WSInit()
	router.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	router.Get("/", websocket.New(ws.WSController))
}
