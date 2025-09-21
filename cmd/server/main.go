package main

import (
	"log"

	"github.com/pharmacy-modernization-project-model/internal/app"
	"github.com/pharmacy-modernization-project-model/internal/platform/config"
)

func main() {
	cfg := config.Load()
	application, err := app.New(cfg)
	if err != nil { log.Fatal(err) }
	if err := application.Run(); err != nil { log.Fatal(err) }
}
