package main

import (
	"log"

	_ "github.com/AppeiYA/consultation-platform/docs"
	app "github.com/AppeiYA/consultation-platform/internal"
)

// @title           Consultation Platform API
// @version         1.0
// @description     API documentation for Consultation Platform.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name session_id

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}