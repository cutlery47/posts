package gql

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func (gh *gqlHandler) getSessionUser(p graphql.ResolveParams) (uuid.UUID, error) {
	seshId, err := idFromArg(p.Args["sesh_id"])
	if err != nil {
		return uuid.UUID{}, err
	}

	return gh.svc.GetSessionUser(p.Context, *seshId)
}

func (gh *gqlHandler) resolveQueryPost(p graphql.ResolveParams) (interface{}, error) {
	id, err := idFromArg(p.Args["id"])
	if err != nil {
		return nil, err
	}

	return gh.svc.GetPost(p.Context, *id)
}

func (gh *gqlHandler) resolveQueryPosts(p graphql.ResolveParams) (interface{}, error) {
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
}

func (gh *gqlHandler) resolveMutationInsertPost(p graphql.ResolveParams) (interface{}, error) {
	userId, err := gh.getSessionUser(p)
	if err != nil {
		return nil, err
	}

	in, err := inPostFromArg(p.Args["in_post"])
	if err != nil {
		return nil, err
	}

	return gh.svc.InsertPost(p.Context, *in, userId)
}

func (gh *gqlHandler) resolveMutationDeletePost(p graphql.ResolveParams) (interface{}, error) {
	userId, err := gh.getSessionUser(p)
	if err != nil {
		return nil, err
	}

	id, err := idFromArg(p.Args["id"])
	if err != nil {
		return nil, err
	}

	return gh.svc.DeletePost(p.Context, *id, userId)
}

func (gh *gqlHandler) resolveMutationUpdatePost(p graphql.ResolveParams) (interface{}, error) {
	userId, err := gh.getSessionUser(p)
	if err != nil {
		return nil, err
	}

	id, err := idFromArg(p.Args["post_id"])
	if err != nil {
		return nil, err
	}

	in, err := inPostFromArg(p.Args["in_post"])
	if err != nil {
		return nil, err
	}

	return gh.svc.UpdatePost(p.Context, *id, userId, *in)
}

func (gh *gqlHandler) resolveMutationInsertComment(p graphql.ResolveParams) (interface{}, error) {
	userId, err := gh.getSessionUser(p)
	if err != nil {
		return nil, err
	}

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

	return gh.svc.InsertComment(p.Context, *postId, userId, parentId, *in)
}

func (gh *gqlHandler) resolveMutationDeleteComment(p graphql.ResolveParams) (interface{}, error) {
	userId, err := gh.getSessionUser(p)
	if err != nil {
		return nil, err
	}

	postId, err := idFromArg(p.Args["post_id"])
	if err != nil {
		return nil, err
	}

	commId, err := idFromArg(p.Args["comm_id"])
	if err != nil {
		return nil, err
	}

	return gh.svc.DeleteComment(p.Context, *postId, *commId, userId)
}

func (gh *gqlHandler) resolveMutationUpdateComment(p graphql.ResolveParams) (interface{}, error) {
	userId, err := gh.getSessionUser(p)
	if err != nil {
		return nil, err
	}

	postId, err := idFromArg(p.Args["post_id"])
	if err != nil {
		return nil, err
	}

	commId, err := idFromArg(p.Args["comm_id"])
	if err != nil {
		return nil, err
	}

	comm, err := inCommentFromArg(p.Args["in_comm"])
	if err != nil {
		return nil, err
	}

	return gh.svc.UpdateComment(p.Context, *postId, *commId, userId, *comm)
}
