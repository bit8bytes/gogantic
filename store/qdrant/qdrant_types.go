package qdrant

type ID = string
type Map = map[string]interface{}

type UpsertPointsRequest struct {
	Batch struct {
		IDs      []string    `json:"ids"`
		Payloads []Map       `json:"payloads,omitempty"`
		Vectors  [][]float32 `json:"vectors"`
	} `json:"batch"`
}

type UpsertPointsResponse struct {
	Time   float32 `json:"time"`
	Status string  `json:"status"`
	Result struct {
		Status      string `json:"status"`
		OperationId int    `json:"operation_id"`
	} `json:"result"`
}

type DeletePointsRequest struct {
	Points []ID
}

type DeletePointsResponse = UpsertPointsResponse

type CountPointsRequest struct{}

type CountPointsResponse struct {
	Time   float32 `json:"time"`
	Status string  `json:"status"`
	Result struct {
		Count int `json:"count"`
	} `json:"result"`
}

type SearchPointsRequest struct {
	Limit          int       `json:"limit"`
	WithPayload    bool      `json:"with_payload"`
	WithVector     bool      `json:"with_vector"`
	Vector         []float32 `json:"vector"`
	ScoreThreshold float32   `json:"score_threshold"`
}

type SearchPointsResponse struct {
	Result []struct {
		Id      ID      `json:"id"`
		Score   float32 `json:"score"`
		Payload Map     `json:"payload"`
	} `json:"result"`
	Status string  `json:"status"`
	Time   float32 `json:"time"`
}

type RetrievePointsRequest struct {
	Ids []ID `json:"ids"`
}

type RetrievePointsResponse struct {
	Time   float32 `json:"time"`
	Status string  `json:"status"`
	Result []struct {
		Id ID `json:"id"`
	} `json:"result"`
}

type CreateCollectionRequest struct {
	Vectors struct {
		Size     int    `json:"size"`
		Distance string `json:"distance"`
	} `json:"vectors"`
	SparseVectors struct{} `json:"sparse_vectors"`
}

type CreateCollectionResponse struct {
	Time   float32 `json:"time"`
	Status string  `json:"status"`
	Result bool    `json:"result"`
}

type CheckCollectionHealthResponse struct {
	Time   float32 `json:"time"`
	Status string  `json:"status"`
}

type DeleteCollectionResponse = CreateCollectionRequest
