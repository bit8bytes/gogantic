package qdrant

import (
	"context"
	"net/http"
)

func (qs *QdrantStore) CreateCollection(ctx context.Context, collection string, payload CreateCollectionRequest) (*CreateCollectionResponse, error) {
	url := qs.Url.String() + "/collections/" + collection
	var response CreateCollectionResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodPut, url, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (qs *QdrantStore) DeleteCollection(ctx context.Context) (*DeleteCollectionResponse, error) {
	url := qs.Url.String() + "/collections/" + qs.Collection
	var response DeleteCollectionResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodDelete, url, "", &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (qs *QdrantStore) GetCollectionHealth(ctx context.Context) error {
	url := qs.Url.String() + "/collections/" + qs.Collection
	var response CheckCollectionHealthResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodGet, url, "", &response)
	if err != nil {
		return err
	}
	return nil
}
