package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bit8bytes/gogantic/llms"
)

type client struct {
	Model Model
}

func New(model Model) *client {
	return &client{Model: model}
}

func (oc *client) Generate(ctx context.Context, messages []llms.Message) (*llms.ContentResponse, error) {
	client := &http.Client{
		Timeout: 240 * time.Second,
	}

	var endpoint = oc.Model.Endpoint
	if oc.Model.Endpoint == "" {
		endpoint = "http://localhost:11434/api/chat"
	}

	request := Chat{
		Model:     oc.Model.Model,
		Messages:  messages,
		Options:   oc.Model.Options,
		Stream:    oc.Model.Stream,
		Format:    oc.Model.Format,
		KeepAlive: oc.Model.KeepAlive,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("error marshaling request")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, errors.New("create request failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("HTTP request failed")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var chatResponse ChatResponse
	var finalResponse ChatResponse
	for {
		if err := decoder.Decode(&chatResponse); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.New("error decoding response")
		}

		finalResponse.Message.Content += chatResponse.Message.Content
		for _, stopSeq := range oc.Model.Stop {
			if strings.Contains(finalResponse.Message.Content, stopSeq) {
				finalResponse.Message.Content = strings.Split(finalResponse.Message.Content, stopSeq)[0]
				finalResponse.Done = true
				return &llms.ContentResponse{Result: finalResponse.Message.Content}, nil
			}
		}
	}

	if chatResponse.Done {
		return &llms.ContentResponse{Result: finalResponse.Message.Content}, nil
	}

	return &llms.ContentResponse{Result: finalResponse.Message.Content}, nil
}

func (oc *client) StreamContent(ctx context.Context, messages []llms.Message, handler llms.StreamHandler) error {
	endpoint := oc.Model.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/chat"
	}

	request := Chat{
		Model:     oc.Model.Model,
		Messages:  messages,
		Options:   oc.Model.Options,
		Stream:    true,
		Format:    oc.Model.Format,
		KeepAlive: oc.Model.KeepAlive,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(string(requestBody)))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)
	var fullContent string

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var chatResponse ChatResponse
		if err := decoder.Decode(&chatResponse); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error decoding response: %w", err)
		}

		fullContent += chatResponse.Message.Content

		shouldStop := false
		for _, stopSeq := range oc.Model.Stop {
			if strings.Contains(fullContent, stopSeq) {
				parts := strings.Split(fullContent, stopSeq)
				chatResponse.Message.Content = parts[0]
				shouldStop = true
				break
			}
		}

		if err := handler(chatResponse.Message.Content, chatResponse.Done || shouldStop); err != nil {
			return fmt.Errorf("handler error: %w", err)
		}

		if chatResponse.Done || shouldStop {
			return nil
		}
	}
}

func (oc *client) GenerateEmbedding(ctx context.Context, prompt string) (llms.EmbeddingResponse, error) {
	client := &http.Client{
		Timeout: 240 * time.Second,
	}

	var endpoint = oc.Model.Endpoint
	if oc.Model.Endpoint == "" {
		endpoint = "http://localhost:11434/api/embeddings"
	}

	request := Prompt{
		Model:     oc.Model.Model,
		Prompt:    prompt,
		Options:   oc.Model.Options,
		Format:    oc.Model.Format,
		KeepAlive: oc.Model.KeepAlive,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return llms.EmbeddingResponse{}, errors.New("error marshaling request")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return llms.EmbeddingResponse{}, errors.New("create request failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return llms.EmbeddingResponse{}, errors.New("http request failed")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return llms.EmbeddingResponse{}, errors.New("error reading response body")
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var embeddingResponse llms.EmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&embeddingResponse)
	if err != nil {
		return llms.EmbeddingResponse{}, errors.New("error decoding response")
	}

	return llms.EmbeddingResponse{Embedding: embeddingResponse.Embedding}, nil
}
