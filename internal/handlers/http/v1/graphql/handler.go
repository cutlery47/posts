package gql

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	post "github.com/cutlery47/posts/internal/storage/post-storage"
	user "github.com/cutlery47/posts/internal/storage/user-storage"
	errh "github.com/cutlery47/posts/pkg/errhandle"
	"github.com/graphql-go/graphql"
)

type gqlHandler struct {
	us user.Storage
	ps post.Storage

	schema graphql.Schema
}

func New(ps post.Storage, us user.Storage) (*gqlHandler, error) {
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	log.Println(id1, id2, id3)

	return &gqlHandler{
		us:     us,
		ps:     ps,
		schema: schema,
	}, nil
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

	queryStirng, ok := queryField.(string)
	if !ok {
		errh.Handle(errors.New("query can't be converted to string"), w)
		return
	}

	res := graphql.Do(graphql.Params{
		Context:       r.Context(),
		Schema:        gh.schema,
		RequestString: queryStirng,
	})

	json.NewEncoder(w).Encode(res)
}
