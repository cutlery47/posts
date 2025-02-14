package gql

import (
	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func (gh *gqlHandler) initSchema() error {
	var inCommentInput = graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "InCommentInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"user_id": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"content": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		},
	)

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

	var inPostInput = graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "InPostInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"user_id": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"content": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"is_mute": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
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
				"NEWEST": &graphql.EnumValueConfig{
					Value: "newest",
				},
				"OLDEST": &graphql.EnumValueConfig{
					Value: "oldest",
				},
				"UPVOTED": &graphql.EnumValueConfig{
					Value: "upvoted",
				},
				"DOWNVOTED": &graphql.EnumValueConfig{
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
						id, err := idFromArg(p.Args["id"])
						if err != nil {
							return nil, err
						}

						return gh.svc.GetPost(p.Context, *id)
					},
				},
				"posts": &graphql.Field{
					Type: graphql.NewNonNull(
						&graphql.List{
							OfType: postType,
						},
					),
					Description: "get sorted + paginated posts",
					Args: graphql.FieldConfigArgument{
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"offset": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"sort_by": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(sortEnum),
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
							v, _ := limitArg.(int)
							limit = &v
						}

						offsetArg, ok := p.Args["offset"]
						if ok {
							v, _ := offsetArg.(int)
							offset = &v
						}

						return gh.svc.GetPosts(p.Context, limit, offset, sortBy)
					},
				},
			},
		},
	)

	var rootMutation = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"insertPost": &graphql.Field{
					Type: graphql.NewNonNull(postType),
					Args: graphql.FieldConfigArgument{
						"in_post": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(inPostInput),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						in, err := inPostFromArg(p.Args["in_post"])
						if err != nil {
							return nil, err
						}

						return gh.svc.InsertPost(p.Context, *in)
					},
				},
				"deletePost": &graphql.Field{
					Type: graphql.ID,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, err := idFromArg(p.Args["id"])
						if err != nil {
							return nil, err
						}

						return gh.svc.DeletePost(p.Context, *id)
					},
				},
				"updatePost": &graphql.Field{
					Type: graphql.NewNonNull(postType),
					Args: graphql.FieldConfigArgument{
						"post_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"in_post": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(inPostInput),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, err := idFromArg(p.Args["post_id"])
						if err != nil {
							return nil, err
						}

						in, err := inPostFromArg(p.Args["in_post"])
						if err != nil {
							return nil, err
						}

						return gh.svc.UpdatePost(p.Context, *id, *in)
					},
				},
				"insertComment": &graphql.Field{
					Type: graphql.NewNonNull(commentType),
					Args: graphql.FieldConfigArgument{
						"post_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"parent_id": &graphql.ArgumentConfig{
							Type: graphql.ID,
						},
						"in_comment": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(inCommentInput),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						postId, err := idFromArg(p.Args["post_id"])
						if err != nil {
							return nil, err
						}

						var (
							parentId *uuid.UUID
						)

						parentIdArg, ok := p.Args["parent_id"]
						if ok {
							p, err := idFromArg(parentIdArg)
							if err != nil {
								return nil, err
							}

							parentId = p
						}

						in, err := inCommentFromArg(p.Args["in_comment"])
						if err != nil {
							return nil, err
						}

						return gh.svc.InsertComment(p.Context, *postId, parentId, *in)
					},
				},
				"deleteComment": &graphql.Field{
					Type: graphql.ID,
					Args: graphql.FieldConfigArgument{
						"post_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"comm_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						postId, err := idFromArg(p.Args["post_id"])
						if err != nil {
							return nil, err
						}

						commId, err := idFromArg(p.Args["comm_id"])
						if err != nil {
							return nil, err
						}

						return gh.svc.DeleteComment(p.Context, *postId, *commId)
					},
				},
			},
		},
	)

	var schemaConfig = graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return err
	}

	gh.schema = schema

	return nil
}
