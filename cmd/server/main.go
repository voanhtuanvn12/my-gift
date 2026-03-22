// Package main is the entry point of the my-gift service.
//
// @title           My Gift API
// @version         1.0
// @description     RESTful API for the My Gift service.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.email  support@my-gift.local
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @servers.url         http://localhost:8080
// @servers.description Local development server
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT token.
package main

import (
	"log"

	"my-gift/configs"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var app *App
	if cfg.App.Env == "dummy" {
		app, err = InitializeAppDummy(cfg)
	} else {
		app, err = InitializeApp(cfg)
	}
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
