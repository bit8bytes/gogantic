package qdrant

import (
	"context"
	"net/http"
)

func (qs *QdrantStore) upsertPoints(ctx context.Context, payload UpsertPointsRequest) (*UpsertPointsResponse, error) {
	url := qs.Url.String() + "/collections/" + qs.Collection + "/points"
	var response UpsertPointsResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodPut, url, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (qs *QdrantStore) DeletePoints(ctx context.Context, payload DeletePointsRequest) (*DeletePointsResponse, error) {
	url := qs.Url.String() + "/collections/" + qs.Collection + "/points/vectors/delete"
	var response DeletePointsResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodPost, url, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (qs *QdrantStore) CountPoints(ctx context.Context) (*CountPointsResponse, error) {
	url := qs.Url.String() + "/collections/" + qs.Collection + "/points/count"
	var response CountPointsResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodPost, url, CountPointsRequest{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (qs *QdrantStore) RetrievePoints(ctx context.Context, payload RetrievePointsRequest) (*RetrievePointsResponse, error) {
	url := qs.Url.String() + "/collections/" + qs.Collection + "/points"
	var response RetrievePointsResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodPost, url, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
