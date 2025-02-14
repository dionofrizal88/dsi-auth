package main

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	"github.com/dionofrizal88/dsi/auth/interfaces/cmd"
	"github.com/dionofrizal88/dsi/auth/pkg/db"
	"github.com/dionofrizal88/dsi/auth/routes"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// @host localhost:8181
// @schemes http
// main init the auth service.
func main() {
	conf := config.GetConfig("config/config.json")

	dbConnection, errDBConnection := db.NewPostgresDBConnection(conf)
	if errDBConnection != nil {
		log.Fatalf("Failed to connect into postgres db %v", errDBConnection)
	}

	redisConnection, errRedisConnection := db.NewRedisConnection(conf)
	if errRedisConnection != nil {
		log.Fatalf("Failed to connect to redis: %v", errRedisConnection)
	}

	dbService := dao.NewDBService(dbConnection)
	entityRegistry := dao.NewRegistry()

	// Auto Migrate
	errAutoMigrate := entityRegistry.AutoMigrate(dbConnection)
	if errAutoMigrate != nil {
		log.Fatalf("Failed to run auto migration: %v", errAutoMigrate)
	}

	// Init app
	app := cmd.NewCli()
	app.Action = func(c *cli.Context) error {
		// Init Router
		e := routes.
			NewRouter(
				routes.WithConfig(conf),
				routes.WithDBService(dbService),
				routes.WithRedisDB(redisConnection),
			).
			Init()

		e.Logger.Fatal(e.Start(conf.AppPort))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Failed to init CLI: %v", err)
	}
}
