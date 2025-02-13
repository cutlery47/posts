package app

import (
	"fmt"
	"log"
	"os"

	"github.com/cutlery47/posts/config"
	v1 "github.com/cutlery47/posts/internal/handlers/http/v1"
	gql "github.com/cutlery47/posts/internal/handlers/http/v1/graphql"
	post "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/cutlery47/posts/internal/storage/post-storage/mem"
	pgpost "github.com/cutlery47/posts/internal/storage/post-storage/postgres"
	user "github.com/cutlery47/posts/internal/storage/user-storage"
	"github.com/cutlery47/posts/internal/storage/user-storage/mock"
	pguser "github.com/cutlery47/posts/internal/storage/user-storage/postgres"
	"github.com/cutlery47/posts/pkg/httpserver"
)

func Run(conf config.App) error {
	errChan := make(chan error, 1)

	log.Println("[SETUP] setting up post storage...")

	ps, err := getPostStorage(conf.PostStorage, errChan)
	if err != nil {
		return fmt.Errorf("[SETUP ERROR] error when setting up post storage: %v", err)
	}

	log.Println("[SETUP] setting up user storage...")

	us, err := getUserStorage(conf.UserStorage)
	if err != nil {
		return fmt.Errorf("[SETUP ERROR] error when setting up user storage: %v", err)
	}

	log.Println("[SETUP] setting up graphql handler...")

	h, err := gql.New(ps, us)
	if err != nil {
		return fmt.Errorf("[SETUP ERROR] error when seting up graphql handler: %v", err)
	}

	log.Println("[SETUP] setting up http server...")

	srv := httpserver.New(conf.HTTPServer, v1.New(h))

	return srv.Run(errChan)
}

func getPostStorage(conf config.PostStorage, errChan chan<- error) (post.Storage, error) {
	switch conf.Type {
	case "mem":
		fd, err := os.OpenFile(conf.DumpDestination, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			return nil, err
		}

		return mem.NewStorage(conf, fd, fd, errChan)
	case "pg":
		return pgpost.NewStorage()
	default:
		return nil, fmt.Errorf("post storage type undefined. supported types: \"pg\" (not impl), \"mem\"")
	}
}

func getUserStorage(conf config.UserStorage) (user.Storage, error) {
	switch conf.Type {
	case "mock":
		return mock.NewStorage(), nil
	case "pg":
		return pguser.NewStorage()
	default:
		return nil, fmt.Errorf("user storage type undefined. supported types: \"pg\" (not impl), \"mock\"")
	}
}
