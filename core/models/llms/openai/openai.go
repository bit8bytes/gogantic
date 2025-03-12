package openai

import (
	"bytes"
	"context"
	"encoding/json"
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
