package graphql

import (
	"net/http"

	post "github.com/cutlery47/posts/internal/storage/post-storage"
	user "github.com/cutlery47/posts/internal/storage/user-storage"
)

type gqlHandler struct {
	us user.Storage
	ps post.Storage
}

func (gh *gqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func New(ps post.Storage, us user.Storage) *gqlHandler {
	return &gqlHandler{
		us: us,
		ps: ps,
	}
}
