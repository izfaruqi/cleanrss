package routes

import (
	"cleanrss/controllers"

	"github.com/gofiber/fiber/v2"
)

func entryRouter(router fiber.Router) {
	router.Get("/", controllers.EntryGetFromDBByAllProviders)
	router.Get("/refresh", controllers.EntryRefreshDBFromAllProviders)
	router.Get("/provider/:id/refresh", controllers.EntryRefreshDBFromProvider)
	router.Get("/provider/:id", controllers.EntryGetFromDBByProvider)
	router.Get("/search", controllers.EntrySearch)
}
