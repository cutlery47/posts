package main

import (
	"log"

	"github.com/cutlery47/posts/config"
	"github.com/cutlery47/posts/internal/app"
)

func main() {
	conf, err := config.New(".env")
	if err != nil {
		log.Fatalf("error when reading config: %v", err)
	}

	err = app.Run(*conf)
	if err != nil {
		log.Fatalf("runtime error: %v", err)
	}

	log.Println("service shut down gracefully")
}
