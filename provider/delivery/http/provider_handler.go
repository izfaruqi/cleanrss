package http

import (
	"cleanrss/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type providerHttpHandler struct {
	U domain.ProviderUsecase
}

func NewProviderHttpHandler(httpRouter fiber.Router, u domain.ProviderUsecase) {
	handler := providerHttpHandler{U: u}
	httpRouter.Get("/", handler.getAll)
	httpRouter.Get("/:id", handler.getById)
	httpRouter.Post("/", handler.insert)
	httpRouter.Post("/:id", handler.update)
	httpRouter.Delete("/:id", handler.delete)
}

func (h providerHttpHandler) getAll(c *fiber.Ctx) error {
	providers, err := h.U.GetAll()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(providers)
}

func (h providerHttpHandler) getById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	provider, err := h.U.GetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(404).JSON(err.Error())
		}
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(provider)
}

func (h providerHttpHandler) insert(c *fiber.Ctx) error {
	provider := new(domain.Provider)
	err := c.BodyParser(provider)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err = h.U.Insert(provider)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(map[string]int64{"id": provider.Id})
}

func (h providerHttpHandler) update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	provider := new(domain.Provider)
	err = c.BodyParser(provider)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	provider.Id = id
	err = h.U.Update(*provider)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func (h providerHttpHandler) delete(c *fiber.Ctx) error {
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
