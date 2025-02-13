package main

import (
	"log"

	"github.com/cutlery47/posts/config"
)

func main() {
	conf, err := config.New(".env")
	if err != nil {
		log.Fatalf("error when reading config: %v", err)
	}

}
