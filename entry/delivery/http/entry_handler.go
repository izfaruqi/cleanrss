package http

import (
	"cleanrss/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type entryHttpHandler struct {
	U domain.EntryUsecase
}

func NewEntryHttpRouter(httpRouter fiber.Router, u domain.EntryUsecase) {
	handler := entryHttpHandler{U: u}
	httpRouter.Get("/refresh", handler.refreshFromAllProviders)
	httpRouter.Get("/refresh/provider/:id", handler.refreshFromProvider)
	httpRouter.Get("/query", handler.getByQuery)
}

func (h entryHttpHandler) refreshFromAllProviders(c *fiber.Ctx) error {
	err := h.U.TriggerRefreshAll()
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func (h entryHttpHandler) refreshFromProvider(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	err = h.U.TriggerRefresh(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func (h entryHttpHandler) getByQuery(c *fiber.Ctx) error {
	var entries []domain.Entry
	query := c.Query("q", "")
	dateFrom, err := strconv.ParseInt(c.Query("date_from", "-1"), 10, 64)
	dateUntil, err := strconv.ParseInt(c.Query("date_until", "-1"), 10, 64)
	providerId, err := strconv.ParseInt(c.Query("provider_id", "-1"), 10, 64)
	limit, err := strconv.ParseInt(c.Query("limit", "40"), 10, 64)
	offset, err := strconv.ParseInt(c.Query("offset", "0"), 10, 64)
	includeAll, err := strconv.ParseBool(c.Query("include_all", "false"))
	allowRefresh, err := strconv.ParseBool(c.Query("allow_refresh", "true"))

	if err != nil {
		return c.Status(400).JSON("Bad Provider ID format.")
	}
	if query == "" {
		return c.Status(400).JSON("Query cannot be empty.")
	}
	entries, err = h.U.GetByQuery(query, dateFrom, dateUntil, providerId, limit, offset, includeAll)
	if (entries == nil || len(entries) == 0) && providerId != -1 && allowRefresh && offset == 0 {
		err = h.U.TriggerRefresh(providerId)
		if err != nil {
			return c.Status(500).JSON(err.Error())
		}
		entries, err = h.U.GetByQuery(query, dateFrom, dateUntil, providerId, limit, offset, includeAll)
	}
	if err != nil {
		return c.Status(500).JSON("Internal server error")
	}
	return c.Status(200).JSON(entries)
}
