package main

import (
	// "errors"
	// "context"
	"fmt"
	// "log"
	"net/http"
	// "os"

	// "server/db"
	"server/env"
	"server/logging"
	// "server/models"

	"server/routes"
	// "server/kafka"

	"github.com/rs/zerolog/log"
)

type Application struct {
	Config env.Config
	// Models models.Models
}

func (app *Application) Serve() error {
	port := app.Config.PORT

	log.Printf("🚀 Server listening on port %s", port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: routes.Routes(),
	}

	return srv.ListenAndServe()
}

func init() {
	_, _ = env.Load()
	logging.InitLogging()
}

// @title Go Notification Service API
// @version 1.0
// @description This is a sample server for a notification service written in Go.
// @termsOfService http://localhost:6000/terms/

// @contact.name API Support
// @contact.url http://www.localhost:6000/support
// @contact.email sushrit.pk21@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:6000
// @BasePath /
func main() {

	app := Application{
		Config: env.DefaultConfig,
		// Models: models.New(dbConn.DB),
	}
	
	err := app.Serve()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error starting server")
	}
}
