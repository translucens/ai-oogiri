package main

import (
	"context"
	"log"

	"github.com/caarlos0/env"
	"github.com/translucens/oogiri/internal/ai"
	"github.com/translucens/oogiri/internal/database"
	"github.com/translucens/oogiri/internal/webserver"
)

type Config struct {
	DBUsername   string `env:"DB_USERNAME,required"`
	DBPassword   string `env:"DB_PASSWORD,required"`
	DBHost       string `env:"DB_HOST"`
	DBPort       int    `env:"DB_PORT"`
	DBUNIXSocket string `env:"DB_UNIX_SOCKET"`
	DBName       string `env:"DB_NAME,required"`
	ServerPort   int    `env:"PORT,required"`
	ProjectID    string `env:"PROJECT_ID,required"`
}

func main() {

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	if (cfg.DBHost == "" || cfg.DBPort == 0) && cfg.DBUNIXSocket == "" {
		log.Fatalf(
			"Either DB_HOST and DB_PORT or DB_UNIX_SOCKET must be set",
		)
	}

	ctx := context.Background()

	dbClient, err := database.NewClient(ctx, cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBUNIXSocket, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	aiClient, err := ai.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer aiClient.Close()

	server, err := webserver.NewServer(dbClient, aiClient, cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to set up server: %v", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
