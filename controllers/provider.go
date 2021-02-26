package controllers

import (
	"cleanrss/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func ProviderGetAll(c *fiber.Ctx) error {
	providers, err := models.ProviderGetAll()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(providers)
}

func ProviderGetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	provider, err := models.ProviderGetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(404).JSON(err.Error())
		}
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(provider)
}

func ProviderInsert(c *fiber.Ctx) error {
	provider := new(models.Provider)
	err := c.BodyParser(provider)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err = validate.StructExcept(provider, "id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	insertedId, err := models.ProviderInsert(provider)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(map[string]int64{"id": insertedId})
}

func ProviderUpdate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	provider := new(models.Provider)
	err = c.BodyParser(provider)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err = validate.StructExcept(provider, "id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	provider.Id = id
	affected, err := models.ProviderUpdate(provider)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if affected < 1 {
		return c.Status(404).JSON("ID not found.")
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}

func ProviderDelete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	affected, err := models.ProviderDelete(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if affected < 1 {
		return c.Status(404).JSON("ID not found.")
	}
	return c.Status(200).JSON(map[string]bool{"success": true})
}