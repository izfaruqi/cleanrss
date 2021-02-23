package controllers

import (
	"cleanrss/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func CleanerGetAll(c *fiber.Ctx) error {
	cleaners, err := models.CleanerGetAll()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.JSON(cleaners)
}

func CleanerGetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	cleaner, err := models.CleanerGetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(404).JSON(err.Error())
		}
		return c.Status(500).JSON(err.Error())
	}
	
	return c.JSON(cleaner)
}

func CleanerInsert(c *fiber.Ctx) error {
	cleaner := new(models.Cleaner)
	err := c.BodyParser(cleaner)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err = validate.StructExcept(cleaner, "id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	insertedId, err := models.CleanerInsert(cleaner)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(map[string]int64{"id": insertedId})
}

func CleanerUpdate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	if id == 0 {
		return c.Status(403).JSON("Default cleaner cannot be edited.")
	}
	cleaner := new(models.Cleaner)
	err = c.BodyParser(cleaner)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err = validate.StructExcept(cleaner, "id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	cleaner.Id = id
	affected, err := models.CleanerUpdate(cleaner)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if affected < 1 {
		return c.Status(404).JSON("ID not found.")
	}
	return c.JSON(map[string]bool{"success": true})
}

func CleanerDelete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON("ID is invalid.")
	}
	if id == 0 {
		return c.Status(403).JSON("Default cleaner cannot be deleted.")
	}
	affected, err := models.CleanerDelete(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if affected < 1 {
		return c.Status(404).JSON("ID not found.")
	}
	return c.JSON(map[string]bool{"success": true})
}