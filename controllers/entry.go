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

func EntryRefreshDBFromAllProviders(c *fiber.Ctx) error {
	services.RefreshEntriesFromProviders()
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func EntryQuery(c *fiber.Ctx) error {
	var entries *[]models.Entry
	query := c.Query("q", "")
	dateFrom, err := strconv.ParseInt(c.Query("date_from", "-1"), 10, 64)
	dateUntil, err := strconv.ParseInt(c.Query("date_until", "-1"), 10, 64)
	providerId, err := strconv.ParseInt(c.Query("provider_id", "-1"), 10, 64)
	limit, err := strconv.ParseInt(c.Query("limit", "40"), 10, 64)
	offset, err := strconv.ParseInt(c.Query("offset", "0"), 10, 64)
	includeAll, err := strconv.ParseBool(c.Query("include_all", "false"))
	allowRefresh, err := strconv.ParseBool(c.Query("allow_refresh", "true"))

	if err != nil {
		c.Status(400).JSON("Bad Provider ID format.")
	}
	if query == "" {
		c.Status(400).JSON("Query cannot be empty.")
	}
	entries, err = models.EntrySearch(query, dateFrom, dateUntil, providerId, limit, offset, includeAll)
	if (len(*entries) == 0 || entries == nil) && providerId != -1 && allowRefresh && offset == 0 {
		err = models.EntryDBRefreshFromProvider(providerId)
		if err != nil {
			return c.Status(500).JSON(err.Error())
		}
		entries, err = models.EntrySearch(query, dateFrom, dateUntil, providerId, limit, offset, includeAll)
	}
	if err != nil {
		c.Status(500).JSON("Internal server error")
	}
	return c.Status(200).JSON(entries)
}
