package errh

import (
	"log"
	"net/http"
)

func Handle(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}
