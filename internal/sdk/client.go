package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GraphQLRequest models the JSON payload sent to the GraphQL API.
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// HiveClient encapsulates an HTTP client, a GraphQL endpoint, and an API token.
type HiveClient struct {
	client   *http.Client
	endpoint string
	token    string
}

// NewHiveClient creates a new HiveClient instance.
func NewHiveClient(client *http.Client, endpoint, token string) *HiveClient {
	return &HiveClient{
		client:   client,
		endpoint: endpoint,
		token:    token,
	}
}

// Execute sends the provided GraphQL query (with variables) to the endpoint
// and returns the response as a map or an error.
func (hc *HiveClient) Execute(ctx context.Context, query string, variables map[string]any, result any) error {
	// Build the request payload.
	payload := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Marshal payload to JSON.
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a new HTTP POST request.
	req, err := http.NewRequest("POST", hc.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers.
	req.Header.Set("Content-Type", "application/json")
	if hc.token != "" {
		req.Header.Set("Authorization", "Bearer "+hc.token)
	}

	// Execute the request.
	resp, err := hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Ensure we received an OK response.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", resp.Status)
	}

	// Decode the JSON response.
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
