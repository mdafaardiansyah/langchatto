package main

import (
	"fmt"
	"github.com/kooroshh/fiber-boostrap/bootstrap"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"log"
)

// main initializes and starts the Fiber application by creating a new
// application instance using the bootstrap package and then listens on
// the specified host and port defined in the environment variables.
// If the application fails to start, it logs the error and terminates
// the program.
func main() {
	app := bootstrap.NewApplication()
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT", "4000"))))
}

