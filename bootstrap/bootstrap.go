package bootstrap

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/kooroshh/fiber-boostrap/app/ws"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"github.com/kooroshh/fiber-boostrap/pkg/router"
	"go.elastic.co/apm"
)

// NewApplication returns a new Fiber app with the following middleware:
// - recover.New(): to recover from panics
// - logger.New(): to log all requests
// - monitor.New(): to expose metrics at /dashboard
// - ws.ServeWSMessaging(): to serve WebSocket connections at /message/v1/send
// - router.InstallRouter(): to install routes for API and HTTP
func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogFile()

	database.SetupDatabase()
	database.SetupMongoDB()

	apm.DefaultTracer.Service.Name = "langchatto-app"
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	go ws.ServeWSMessaging(app)

	router.InstallRouter(app)

	return app
}

// SetupLogFile configures the logging system to write logs to both the standard
// output and a file named "langchatto-app.log" located in the "logs" directory.
// If the log file does not exist, it will be created. If there is an error
// opening or creating the log file, the function will log a fatal error and
// terminate the program.
func SetupLogFile() {
	logFile, err := os.OpenFile("./logs/langchatto-app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
