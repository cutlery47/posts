package gql

import (
	"time"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var ts = time.Now()

var (
	id1 = uuid.New()
	id2 = uuid.New()
	id3 = uuid.New()

	commid1 = uuid.New()
	commid2 = uuid.New()

	replid1 = uuid.New()
)

var psts = map[uuid.UUID]storage.Post{
	id1: {
		Id: id1,
		Comments: map[uuid.UUID]storage.Comment{
			commid1: {
				Id: commid1,
				Replies: map[uuid.UUID]storage.Comment{
					replid1: {
						Id: replid1,
					},
				},
			},
		},
	},
	id2: {
		Id: id2,
		Comments: map[uuid.UUID]storage.Comment{
			commid2: {
				Id: commid2,
			},
		},
	},
	id3: {
		Id: id3,
	},
}

//=============================

var inCommentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "InComment",
		Fields: graphql.Fields{
			"user_id": &graphql.Field{
				Type: graphql.ID,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var commentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"upvotes": &graphql.Field{
				Type: graphql.Int,
			},
			"downvotes": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"deleted_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var inPostType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "InPost",
		Fields: graphql.Fields{
			"user_id": &graphql.Field{
				Type: graphql.ID,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"is_mute": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

var postType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"in_post": &graphql.Field{
				Type: inPostType,
			},
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"upvotes": &graphql.Field{
				Type: graphql.Int,
			},
			"downvotes": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"deleted_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"post": &graphql.Field{
				Type:        postType,
				Description: "get post by its id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"]

					strId, ok := id.(string)
					if !ok {
						return nil, ErrCantConvertToUUID
					}

					uuidId, err := uuid.Parse(strId)
					if err != nil {
						return nil, ErrCantConvertToUUID
					}

					tmp, ok := psts[uuidId]
					if !ok {
						return nil, nil
					}

					return tmp, nil
				},
			},
			"posts": &graphql.Field{
				Type: &graphql.List{
					OfType: postType,
				},
				Description: "get all posts",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					slicePsts := make([]storage.Post, 0, len(psts))

					for _, v := range psts {
						slicePsts = append(slicePsts, v)
					}

					return slicePsts, nil
				},
			},
		},
	},
)

var schemaConfig = graphql.SchemaConfig{
	Query: rootQuery,
}

// adding cyclic fields
func init() {
	commentType.AddFieldConfig(
		"replies",
		&graphql.Field{
			Type: &graphql.List{
				OfType: commentType,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				src := p.Source.(storage.Comment)

				var (
					repls = make([]storage.Comment, 0, len(src.Replies))
				)

				for _, v := range src.Replies {
					repls = append(repls, v)
				}

				return repls, nil
			},
		})

	postType.AddFieldConfig(
		"comments",
		&graphql.Field{
			Type: &graphql.List{
				OfType: commentType,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				src := p.Source.(storage.Post)

				var (
					comms = make([]storage.Comment, 0, len(src.Comments))
				)

				for _, v := range src.Comments {
					comms = append(comms, v)
				}

				return comms, nil
			},
		},
	)
}
