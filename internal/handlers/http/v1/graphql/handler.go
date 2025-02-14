package gql

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cutlery47/posts/internal/service"
	errh "github.com/cutlery47/posts/pkg/errhandle"
	"github.com/graphql-go/graphql"
)

type gqlHandler struct {
	svc *service.Service

	schema graphql.Schema
}

func New(svc *service.Service) (*gqlHandler, error) {
	gh := &gqlHandler{
		svc: svc,
	}

	if err := gh.initSchema(); err != nil {
		return nil, err
	}

	return gh, nil
}

func (gh *gqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}

	queryJson := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&queryJson)
	if err != nil {
		errh.Handle(err, w)
		return
	}

	queryField, ok := queryJson["query"]
	if !ok {
		errh.Handle(errors.New("query was not provided"), w)
		return
	}

	queryString, ok := queryField.(string)
	if !ok {
		errh.Handle(errors.New("query can't be converted to string"), w)
		return
	}

	res := graphql.Do(graphql.Params{
		Context:       r.Context(),
		Schema:        gh.schema,
		RequestString: queryString,
	})

	json.NewEncoder(w).Encode(res)
}
