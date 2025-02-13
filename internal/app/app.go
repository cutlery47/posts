package app

import (
	"net/http"

	"github.com/cutlery47/posts/config"
	"github.com/cutlery47/posts/pkg/httpserver"
)

func Run(conf config.App) error {
	errChan := make(chan error, 1)

	srv := httpserver.New(conf.HTTPServer, &http.ServeMux{})

	return srv.Run(errChan)
}
