package controllers

import "github.com/gofiber/fiber/v2"

// RenderUI renders the "index" HTML page using the fiber context.
// It does not pass any data to the template.
// Returns an error if the rendering process fails.
func RenderUI(c *fiber.Ctx) error {
	return c.Render("index", nil)
}
