package qdrant

import (
	"context"
	"net/http"
)

func (qs *QdrantStore) searchPoints(ctx context.Context, payload SearchPointsRequest) (*SearchPointsResponse, error) {
	url := qs.Url.String() + "/collections/" + qs.Collection + "/points/search"
	var response SearchPointsResponse
	err := qs.sendHttpRequestWithContext(ctx, http.MethodPost, url, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
