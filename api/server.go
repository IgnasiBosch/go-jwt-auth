package api

import (
	"github.com/IgnasiBosch/go-jwt-auth/api/controllers"
	"github.com/IgnasiBosch/go-jwt-auth/api/seed"
	"os"
)

var server = controllers.Server{}

func Run() {
	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	if os.Getenv("DB_RUN_SEED") == "true" {
		seed.Load(server.DB)
	}

	server.Run(":9090")
}
