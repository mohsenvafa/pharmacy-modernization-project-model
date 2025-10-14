package main

import (
	"fmt"
	"log"
	"os"

	"pharmacy-modernization-project-model/internal/app"
	"pharmacy-modernization-project-model/internal/platform/config"
)

func main() {
	fmt.Printf("CurrentProcess ID: %d\n", os.Getpid())
	cfg := config.Load()
	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
