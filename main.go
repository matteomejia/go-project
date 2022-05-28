package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Workshop struct {
	gorm.Model
	Name    string  `json:"name"`
	Level   int     `json:"level"`
	Start   string  `json:"start"`
	Price   float32 `json:"price"`
	Private bool    `json:"private" gorm:"default:false"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Workshop{})

	app := fiber.New(fiber.Config{
		Prefork:      false,
		UnescapePath: true,
		AppName:      "CreaCode",
	})

	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/workshops", func(c *fiber.Ctx) error {
		var workshops []Workshop
		db.Find(&workshops)
		return c.Status(200).JSON(workshops)
	})

	app.Post("/workshops", func(c *fiber.Ctx) error {
		workshop := new(Workshop)
		if err := c.BodyParser(workshop); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		db.Create(&workshop)
		return c.Status(201).JSON(workshop)
	})

	app.Get("/workshops/:id", func(c *fiber.Ctx) error {
		var workshop Workshop
		id, idErr := c.ParamsInt("id")
		if idErr != nil {
			return c.Status(400).SendString("Bad Request")
		}
		if err := db.Find(&workshop, id).Error; err != nil {
			return c.Status(404).SendString(err.Error())
		}
		return c.Status(200).JSON(workshop)
	})

	app.Listen(":3000")
}
