package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of HTTPClient for testing.
type mockHTTPClient struct{}

const mockData = `{
	"items": [
		{
			"html_url": "https://github.com/test/repo/blob/main/file.go",
			"repository": {
				"full_name": "test/repo"
			}
		}
	]
}`

func (m *mockHTTPClient) DoRequest(ctx context.Context, client *http.Client, url string, headers map[string]string, body interface{}, response interface{}) error {
	return json.Unmarshal([]byte(mockData), response)
}

// Test GitHubClient's SearchCode function.
func TestSearchCode(t *testing.T) {
	// Start a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	// Set up test environment
	client := NewGitHubClient(WithBaseURL(mockServer.URL), WithToken("github-search-test-token"))
	// Call the function
	data, err := client.SearchCode(context.Background(), "test-query")
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data.Items, 1)
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.JSONEq(t, mockData, string(jsonData))
}

// Test SearchCode with API error response
func TestSearchCodeAPIError(t *testing.T) {
	// Start a mock HTTP server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "Internal Server Error"}`))
	}))
	defer mockServer.Close()

	// Set up test environment
	maximumRetryDelay := 1 * time.Second
	client := NewGitHubClient(WithBaseURL(mockServer.URL), WithToken("github-search-token"), WithMaximumRetryDelay(&maximumRetryDelay))
	// Call the function
	data, err := client.SearchCode(context.Background(), "test-query")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, data)
}
