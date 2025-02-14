package gql

import (
	"encoding/json"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
)

func idFromArg(arg any) (*uuid.UUID, error) {
	idStr, ok := arg.(string)
	if !ok {
		return nil, ErrCantConvertToUUID
	}

	v, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return &v, err
}

func inCommentFromArg(arg any) (*storage.InComment, error) {
	argJson, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}

	var in storage.InComment

	err = json.Unmarshal(argJson, &in)
	if err != nil {
		return nil, err
	}

	return &in, err
}

func inPostFromArg(arg any) (*storage.InPost, error) {
	argJson, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}

	var in storage.InPost

	err = json.Unmarshal(argJson, &in)
	if err != nil {
		return nil, err
	}

	return &in, err
}
