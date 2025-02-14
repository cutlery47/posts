package gql

import (
	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// var ts = time.Now()

// var (
// 	id1 = uuid.New()
// 	id2 = uuid.New()
// 	id3 = uuid.New()

// 	commid1 = uuid.New()
// 	commid2 = uuid.New()

// 	replid1 = uuid.New()
// )

// var psts = map[uuid.UUID]storage.Post{
// 	id1: {
// 		Id: id1,
// 		InPost: storage.InPost{
// 			Content: "post1",
// 		},
// 		Comments: map[uuid.UUID]storage.Comment{
// 			commid1: {
// 				Id: commid1,
// 				InComment: storage.InComment{
// 					Content: "comment1",
// 				},
// 				Replies: map[uuid.UUID]storage.Comment{
// 					replid1: {
// 						Id: replid1,
// 						InComment: storage.InComment{
// 							Content: "comment2",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// 	id2: {
// 		Id: id2,
// 		Comments: map[uuid.UUID]storage.Comment{
// 			commid2: {
// 				Id: commid2,
// 			},
// 		},
// 	},
// 	id3: {
// 		Id: id3,
// 	},
// }

func (gh *gqlHandler) initSchema() error {
	var inCommentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "InComment",
			Fields: graphql.Fields{
				"user_id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"content": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		},
	)

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"in_comment": &graphql.Field{
					Type: graphql.NewNonNull(inCommentType),
				},
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"upvotes": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"downvotes": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"created_at": &graphql.Field{
					Type: graphql.NewNonNull(graphql.DateTime),
				},
				"updated_at": &graphql.Field{
					Type: graphql.NewNonNull(graphql.DateTime),
				},
				"deleted_at": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)

	commentType.AddFieldConfig(
		"replies",
		&graphql.Field{
			Type: &graphql.List{
				OfType: commentType,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// convert reply map into slice
				src := p.Source.(storage.Comment)

				var (
					repls = make([]storage.Comment, 0, len(src.Replies))
				)

				for _, v := range src.Replies {
					repls = append(repls, v)
				}

				return repls, nil
			},
		},
	)

	var inPostType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "InPost",
			Fields: graphql.Fields{
				"user_id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"content": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"is_mute": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
			},
		},
	)

	var postType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Post",
			Fields: graphql.Fields{
				"in_post": &graphql.Field{
					Type: graphql.NewNonNull(inPostType),
				},
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"upvotes": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"downvotes": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"created_at": &graphql.Field{
					Type: graphql.NewNonNull(graphql.DateTime),
				},
				"updated_at": &graphql.Field{
					Type: graphql.NewNonNull(graphql.DateTime),
				},
				"deleted_at": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)

	postType.AddFieldConfig(
		"comments",
		&graphql.Field{
			Type: &graphql.List{
				OfType: commentType,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// convert comment map into slice
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

	var sortEnum = graphql.NewEnum(
		graphql.EnumConfig{
			Name: "SortEnum",
			Values: graphql.EnumValueConfigMap{
				"newest": &graphql.EnumValueConfig{
					Value: "newest",
				},
				"oldest": &graphql.EnumValueConfig{
					Value: "oldest",
				},
				"upvoted": &graphql.EnumValueConfig{
					Value: "upvoted",
				},
				"downvoted": &graphql.EnumValueConfig{
					Value: "downvoted",
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
						idArg, _ := p.Args["id"]

						idStr, ok := idArg.(string)
						if !ok {
							return nil, ErrCantConvertToUUID
						}

						id, err := uuid.Parse(idStr)
						if err != nil {
							return nil, ErrCantConvertToUUID
						}

						return gh.svc.GetPost(p.Context, id)
					},
				},
				"posts": &graphql.Field{
					Type: graphql.NewNonNull(
						&graphql.List{
							OfType: postType,
						},
					),
					Description: "get all posts",
					Args: graphql.FieldConfigArgument{
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"offset": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"sort_by": &graphql.ArgumentConfig{
							Type: sortEnum,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var (
							limit  *int
							offset *int
							sortBy string
						)

						sortBy = p.Args["sort_by"].(string)

						limitArg, ok := p.Args["limit"]
						if ok {
							limit = limitArg.(*int)
						}

						offsetArg, ok := p.Args["offset"]
						if !ok {
							offset = offsetArg.(*int)
						}

						return gh.svc.ListPosts(p.Context, limit, offset, sortBy)
					},
				},
			},
		},
	)

	var schemaConfig = graphql.SchemaConfig{
		Query: rootQuery,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return err
	}

	gh.schema = schema

	return nil
}
