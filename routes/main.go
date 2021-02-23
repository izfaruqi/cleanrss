package routes

import (
	"cleanrss/controllers"

	"github.com/gofiber/fiber/v2"
)


func RoutesInit(router fiber.Router){
	controllers.ValidatorInit()
	router.Get("/", func(c *fiber.Ctx) error {
		c.Status(200)
		return c.JSON(map[string]string{"version": "0.1.0"})
	})
	router.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})
	
	providerGroup := router.Group("/provider")
	providerRouter(providerGroup)
	cleanerGroup := router.Group("/cleaner")
	cleanerRouter(cleanerGroup)
	entryGroup := router.Group("/entry")
	entryRouter(entryGroup)
}