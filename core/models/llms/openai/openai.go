package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/bit8bytes/gogantic/core/models"
)

type OpenAiClient struct {
	Model Model
}

func New(model Model) *OpenAiClient {
	return &OpenAiClient{
		Model: model,
	}
}

func (oc *OpenAiClient) GenerateContent(ctx context.Context, messages []models.MessageContent) (models.ContentResponse, error) {
	httpClient := &http.Client{
		Timeout: 240 * time.Second,
	}

	requestPayload := OpenAIRequest{
		Model:       oc.Model.Model,
		Messages:    messages,
		Temperature: oc.Model.Temperature,
		Stream:      oc.Model.Stream,
		Stop:        oc.Model.Stop,
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return models.ContentResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return models.ContentResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oc.Model.APIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return models.ContentResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ContentResponse{}, err
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return models.ContentResponse{}, err
	}

	if len(openAIResp.Choices) == 0 {
		return models.ContentResponse{}, err
	}

	contentResponse := models.ContentResponse{
		Result: openAIResp.Choices[0].Message.Content,
	}

	return contentResponse, nil
}

func (oc *OpenAiClient) GenerateEmbedding(ctx context.Context, input string) (models.EmbeddingResponse, error) {
	client := &http.Client{
		Timeout: 240 * time.Second,
	}

	endpoint := "https://api.openai.com/v1/embeddings"

	requestPayload := EmbeddingRequest{
		Model:          oc.Model.Model,
		Input:          input,
		EncodingFormat: "float",
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return models.EmbeddingResponse{}, errors.New("error marshaling request")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return models.EmbeddingResponse{}, errors.New("create request failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oc.Model.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return models.EmbeddingResponse{}, errors.New("http request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return models.EmbeddingResponse{}, errors.New("API error: " + string(bodyBytes))
	}

	var embeddingResponse EmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&embeddingResponse)
	if err != nil {
		return models.EmbeddingResponse{}, errors.New("error decoding response")
	}

	// Extract the first embedding from the response
	if len(embeddingResponse.Data) == 0 {
		return models.EmbeddingResponse{}, errors.New("no embedding returned")
	}

	// Assuming your models.EmbeddingResponse just needs the embedding vector
	return models.EmbeddingResponse{Embedding: embeddingResponse.Data[0].Embedding}, nil
}
