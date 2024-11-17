package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/kooroshh/fiber-boostrap/app/controllers"
)

type HttpRouter struct {
}

// InstallRouter registers a single route under / that renders the index.html template.
// It also enables CORS and CSRF protection on the route.
func (h HttpRouter) InstallRouter(app *fiber.App) {
	group := app.Group("", cors.New(), csrf.New())
	group.Get("/", controllers.RenderUI)
}

// NewHttpRouter creates and returns a new instance of HttpRouter.
// The HttpRouter is responsible for setting up HTTP routes and middleware such as CORS and CSRF protection.
func NewHttpRouter() *HttpRouter {
	return &HttpRouter{}
}
