package router

import (
	"github.com/gofiber/fiber/v2"
)

// InstallRouter initializes the application by setting up API and HTTP routers.
// It configures all necessary routes and middleware for handling HTTP requests.
func InstallRouter(app *fiber.App) {
	setup(app, NewApiRouter(), NewHttpRouter())
}

// setup takes a Fiber app and one or more Routers as arguments.
// It then calls InstallRouter on each of the routers, passing the app as an argument.
// This is a convenience function for setting up the application with multiple routers.
func setup(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
