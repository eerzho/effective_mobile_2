package main

import (
	"log"

	"effective_mobile_2/internal/app/http_srv"
	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/config"
	"effective_mobile_2/internal/database"
)

func main() {
	log.Print("parsing configuration")
	if err := config.Parse(); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	log.Print("connecting database")
	if err := database.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Print("migrating")
	if err := database.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Print("set upping logger")
	app_log.Setup(config.Cfg().Logger.Level)

	log.Print("running http server")
	if err := http_srv.Run(); err != nil {
		log.Fatalf("stoping http server: %v", err)
	}
	log.Print("stopped http server")
}
