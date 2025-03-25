package httpClient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Thesohan/weaveGitHubSearchService/server/constants"
	"github.com/cenkalti/backoff/v4"
)

// HTTPClient defines an interface for making HTTP requests.
type HTTPClient interface {
	DoRequest(ctx context.Context, client *http.Client, url string, headers map[string]string, body interface{}, response interface{}) error
}

// GetRequest struct for handling GET requests.
type GetRequest struct{}

// DoRequest executes an HTTP GET request.
func (g *GetRequest) DoRequest(ctx context.Context, client *http.Client, url string, headers map[string]string, body interface{}, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, constants.HTTP_METHOD_GET, url, nil)
	if err != nil {
		return err
	}
	setHeaders(req, headers)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = constants.MAXIUM_RETRY_DELAY // Set max retry duration

	operation := func() error {
		fmt.Printf("github api call at time: %v\n", time.Now().Format(time.DateTime))
		err := executeRequest(client, req, response)
		if err != nil {
			fmt.Printf("error while calling search api: %v\n", err)
		}
		return err
	}

	return backoff.Retry(operation, bo)
}

// Helper function to set headers for a request.
func setHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// Helper function to execute an HTTP request and parse response.
func executeRequest(client *http.Client, req *http.Request, response interface{}) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}

// Factory function to create appropriate HTTP request type.
func NewHTTPRequest(method string) (HTTPClient, error) {
	switch method {
	case constants.HTTP_METHOD_GET:
		return &GetRequest{}, nil
	// Add cases for PUT, DELETE, etc.
	default:
		return nil, fmt.Errorf("Not defined")
	}
}
