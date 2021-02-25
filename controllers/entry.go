package controllers

import (
	"cleanrss/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func EntryRefreshDBFromProvider(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	err = models.EntryDBRefreshFromProvider(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func EntryGetFromDBByProvider(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	includeRawJson, err := strconv.ParseBool(c.Query("include_raw_json", "false"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	allowRefresh, err := strconv.ParseBool(c.Query("allow_refresh", "true"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	
	var entries *[]models.Entry
	entries, err = models.EntryGetFromDB(id, limit, offset, includeRawJson)
	if (len(*entries) == 0 || entries == nil) && allowRefresh {
		err = models.EntryDBRefreshFromProvider(id)
		if err != nil {
			return c.Status(500).JSON(err.Error())
		}
		entries, err = models.EntryGetFromDB(id, limit, offset, includeRawJson)
	}
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(entries)
}

func EntryGetFromDBByAllProviders(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	includeRawJson, err := strconv.ParseBool(c.Query("include_raw_json", "false"))
	if err != nil {
		return c.Status(400).JSON("Query params is invalid.")
	}
	
	var entries *[]models.Entry
	entries, err = models.EntryGetFromDB(-1, limit, offset, includeRawJson)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(entries)
}