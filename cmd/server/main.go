package main

import (
	"log"

	"pharmacy-modernization-project-model/internal/app"
	"pharmacy-modernization-project-model/internal/platform/config"
)

func main() {
	cfg := config.Load()
	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
