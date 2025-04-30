package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bit8bytes/gogantic/core/models"
)

// OpenAiClient represents a client for interacting with the OpenAI API
type OpenAiClient struct {
	Model      Model
	HttpClient *http.Client
}

// OpenAiClientOption defines a function type for configuring OpenAiClient
type OpenAiClientOption func(*OpenAiClient)

// WithCustomHttpClient allows setting a custom HTTP client
func WithCustomHttpClient(httpClient *http.Client) OpenAiClientOption {
	return func(c *OpenAiClient) {
		c.HttpClient = httpClient
	}
}

// New creates a new OpenAI client with the specified model and options
func New(model Model, opts ...OpenAiClientOption) *OpenAiClient {
	// Create client with default settings
	client := &OpenAiClient{
		Model: model,
		HttpClient: &http.Client{
			Timeout: 240 * time.Second,
		},
	}

	// Apply any provided options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// GenerateContent sends a request to generate content based on the provided messages
func (oc *OpenAiClient) GenerateContent(ctx context.Context, messages []models.MessageContent) (*models.ContentResponse, error) {
	requestPayload := OpenAIRequest{
		Model:    oc.Model.Model,
		Messages: messages,
	}

	if oc.Model.Temperature != nil {
		requestPayload.Temperature = *oc.Model.Temperature
	}

	if oc.Model.Stream != nil {
		requestPayload.Stream = *oc.Model.Stream
	}

	if oc.Model.Stop != nil {
		requestPayload.Stop = *oc.Model.Stop
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oc.Model.APIKey)

	resp, err := oc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New("API error: " + string(bodyBytes))
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) == 0 {
		return nil, errors.New("no choices returned from API")
	}

	contentResponse := models.ContentResponse{
		Result: openAIResp.Choices[0].Message.Content,
	}

	return &contentResponse, nil
}

func (oc *OpenAiClient) StreamContent(ctx context.Context, messages []models.MessageContent, streamHandler models.StreamHandler) error {
	return fmt.Errorf("streaming not implemented")
}

// GenerateEmbedding sends a request to generate an embedding for the provided input
func (oc *OpenAiClient) GenerateEmbedding(ctx context.Context, input string) (models.EmbeddingResponse, error) {
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

	resp, err := oc.HttpClient.Do(req)
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

	return models.EmbeddingResponse{Embedding: embeddingResponse.Data[0].Embedding}, nil
}
