package routes

import (
	"cleanrss/controllers"

	"github.com/gofiber/fiber/v2"
)

func cleanerRouter(router fiber.Router){
	router.Post("/", controllers.CleanerInsert)
	
	router.Get("/", controllers.CleanerGetAll)

	router.Get("/:id", controllers.CleanerGetById)

	router.Post("/:id", controllers.CleanerUpdate)
	
	router.Delete("/:id", controllers.CleanerDelete)
}