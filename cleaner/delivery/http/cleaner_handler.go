package http

import (
	"cleanrss/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type cleanerHttpHandler struct {
	U domain.CleanerUsecase
}

func NewCleanerHttpHandler(httpRouter fiber.Router, u domain.CleanerUsecase) {
	handler := cleanerHttpHandler{U: u}
	httpRouter.Get("/", handler.getAll)
	httpRouter.Get("/:id", handler.getById)
	httpRouter.Post("/", handler.insert)
	httpRouter.Post("/:id", handler.update)
	httpRouter.Delete("/:id", handler.delete)
	httpRouter.Get("/clean/:id", handler.cleanPage)
}

func (h cleanerHttpHandler) cleanPage(c *fiber.Ctx) error {
	entryId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	cleaner, err := h.U.GetCleanedEntry(entryId)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	return c.Status(200).SendString(cleaner)
}
func (h cleanerHttpHandler) getAll(c *fiber.Ctx) error {
	cleaners, err := h.U.GetAll()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(cleaners)
}

func (h cleanerHttpHandler) getById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	cleaner, err := h.U.GetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(404).JSON(err.Error())
		}
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(cleaner)
}

func (h cleanerHttpHandler) insert(c *fiber.Ctx) error {
	cleaner := new(domain.Cleaner)
	err := c.BodyParser(cleaner)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err = h.U.Insert(cleaner)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(map[string]int64{"id": cleaner.Id})
}

func (h cleanerHttpHandler) update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	cleaner := new(domain.Cleaner)
	err = c.BodyParser(cleaner)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	cleaner.Id = id
	err = h.U.Update(*cleaner)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func (h cleanerHttpHandler) delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	err = h.U.Delete(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}
