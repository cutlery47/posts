package gql

import (
	storage "github.com/cutlery47/posts/internal/storage/post-storage"
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
			Type: graphql.NewNonNull(&graphql.List{
				OfType: commentType,
			}),
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
			Type: graphql.NewNonNull(&graphql.List{
				OfType: commentType,
			}),
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
					Resolve: gh.resolveQueryPost,
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
					Resolve: gh.resolveQueryPosts,
				},
			},
		},
	)

	var seshToken = &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.ID),
	}

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
						"sesh_id": seshToken,
					},
					Resolve: gh.resolveMutationInsertPost,
				},
				"deletePost": &graphql.Field{
					Type: graphql.ID,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"sesh_id": seshToken,
					},
					Resolve: gh.resolveMutationDeletePost,
				},
				"updatePost": &graphql.Field{
					Type: postType,
					Args: graphql.FieldConfigArgument{
						"post_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"in_post": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(inPostInput),
						},
						"sesh_id": seshToken,
					},
					Resolve: gh.resolveMutationUpdatePost,
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
						"sesh_id": seshToken,
					},
					Resolve: gh.resolveMutationInsertComment,
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
						"sesh_id": seshToken,
					},
					Resolve: gh.resolveMutationDeleteComment,
				},
				"updateComment": &graphql.Field{
					Type: commentType,
					Args: graphql.FieldConfigArgument{
						"post_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"comm_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
						"in_comm": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(inCommentInput),
						},
						"sesh_id": seshToken,
					},
					Resolve: gh.resolveMutationUpdateComment,
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
