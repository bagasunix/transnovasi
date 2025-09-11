package main

import (
	"github.com/bagasunix/transnovasi/internal/configs"
)

// @title Transnovasi API
// @version 1.0
// @description API untuk sistem Transnovasi
// @termsOfService http://swagger.io/terms/

// @contact.name Developer Support
// @contact.email support@transnovasi.local

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	configs.Run()
}
