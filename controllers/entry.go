package controllers

import (
	"cleanrss/models"
	"cleanrss/services"
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
	if (len(*entries) == 0 || entries == nil) && allowRefresh && offset == 0 {
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

func EntryRefreshDBFromAllProviders(c *fiber.Ctx) error {
	services.RefreshEntriesFromProviders()
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func EntrySearch(c *fiber.Ctx) error {
	query := c.Query("q", "")
	dateFrom, err := strconv.ParseInt(c.Query("date_from", "-1"), 10, 64)
	dateUntil, err := strconv.ParseInt(c.Query("date_until", "-1"), 10, 64)
	providerId, err := strconv.ParseInt(c.Query("provider_id", "-1"), 10, 64)
	if err != nil {
		c.Status(400).JSON("Bad Provider ID format.")
	}
	if query == "" {
		c.Status(400).JSON("Query cannot be empty.")
	}
	entries, err := models.EntrySearch(query, dateFrom, dateUntil, providerId)
	if err != nil {
		c.Status(500).JSON("Internal server error")
	}
	return c.Status(200).JSON(entries)
}
