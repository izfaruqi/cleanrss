package routes

import (
	"cleanrss/controllers"

	"github.com/gofiber/fiber/v2"
)

func entryRouter(router fiber.Router){
	router.Get("/provider/:id/refresh", controllers.EntryRefreshDBFromProvider)
	router.Get("/provider/:id", controllers.EntryGetFromDBByProvider)
}