package routes

import (
	"cleanrss/controllers"

	"github.com/gofiber/fiber/v2"
)

func providerRouter(router fiber.Router){
	router.Post("/", controllers.ProviderInsert)
	
	router.Get("/", controllers.ProviderGetAll)

	router.Get("/:id", controllers.ProviderGetById)

	router.Post("/:id", controllers.ProviderUpdate)
	
	router.Delete("/:id", controllers.ProviderDelete)
}