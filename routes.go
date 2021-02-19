package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RoutesInit(server *fiber.App) {
	server.Get("/", func(c *fiber.Ctx) error {
		versionInfo := make(map[string]string)
		versionInfo["version"] = "0.1.0"
		return c.JSON(versionInfo)
	})

	server.Get("/shutdown", func(c *fiber.Ctx) error { 
		defer func (){
			_ = server.Shutdown()
		}()
		return c.SendStatus(200)
	})

	providerRoutes(server)
	cleanerRoutes(server)
}

func providerRoutes(server *fiber.App){
	// Get all providers
	server.Post("/provider", func(c *fiber.Ctx) error {
		provider := new(Provider)
		err := c.BodyParser(provider)
		if err != nil {
			return ErrorResponseFactory(400, "JSON_INVALID", err, c)
		}
		id, err := ProviderInsert(provider)
		if err != nil {
			return ErrorResponseFactory(500, "SQL_ERROR", err, c)
		}
		return c.JSON(map[string]int64 { "id": id })
	})
	
	server.Get("/provider", func(c *fiber.Ctx) error {
		providers, err := ProviderGetAll()
		if err != nil {
			return ErrorResponseFactory(500, "INTERNAL_ERROR", err, c)
		}
		if(len(providers) == 0){
			return c.JSON(make([]int, 0))
		} else {
			return c.JSON(providers)
		}
	})

	server.Get("/provider/:id", func(c *fiber.Ctx) error {
		idInt64, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return ErrorResponseFactory(400, "MALFORMED_REQUEST", err, c)
		}
		provider, err := ProviderGetById(idInt64)
		if err != nil {
			return ErrorResponseFactory(500, "INTERNAL_ERROR", err, c)
		}
		return c.JSON(provider)
	})

	server.Get("/provider/:id/entries", func(c *fiber.Ctx) error {
		idInt64, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return ErrorResponseFactory(400, "MALFORMED_REQUEST", err, c)
		}

		var entries *[]Entry
		if(c.Query("fresh", "0") == "1"){
			err = ProviderRefreshEntriesForDB(idInt64)
			if err != nil {
				return ErrorResponseFactory(500, "INTERNAL_ERROR", err, c)
			}
		}

		entries, err = ProviderGetDBEntries(idInt64, 40)
		if err != nil {
			return ErrorResponseFactory(500, "INTERNAL_ERROR", err, c)
		}

		return c.JSON(entries)
	})

	server.Get("/provider/:id/refresh", func(c *fiber.Ctx) error {
		idInt64, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return ErrorResponseFactory(400, "MALFORMED_REQUEST", err, c)
		}
		ProviderRefreshEntriesForDB(idInt64)
		
		return c.JSON(1)
	})
}

func cleanerRoutes(server *fiber.App){
	server.Get("/cleaner/:id", func(c *fiber.Ctx) error {
		idInt64, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return ErrorResponseFactory(400, "MALFORMED_REQUEST", err, c)
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send([]byte(GetCleanPage(idInt64)))
	})

	server.Get("/cleaner", func(c *fiber.Ctx) error {
		parsers, err := GetAllParsers()
		if err != nil {
			return ErrorResponseFactory(404, "NOT_FOUND", err, c)
		}
		return c.JSON(parsers)
	})
}