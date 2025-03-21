package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Thesohan/weaveGitHubSearchService/server/constants"
	httpClient "github.com/Thesohan/weaveGitHubSearchService/server/http_client"
	"github.com/Thesohan/weaveGitHubSearchService/server/structures"
)

// GitHubAPIURL is the base URL for GitHub search API.
const GitHubAPIURL = "https://api.github.com/search/code"

var (
	once      sync.Once
	singleton IGithubClient
)

type githubClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// IGithubClient interface for API requests.
type IGithubClient interface {
	SearchCode(ctx context.Context, query string) (*structures.Data, error)
}

/*
	NewGitHubClient returns a singleton instance of GitHubClient.

baseURL is accepted here to make testing easy, we can pass a mock URL while creating tests
*/
func NewGitHubClient(baseURL ...string) IGithubClient {
	once.Do(func() {
		fmt.Println("Creating new github client")
		url := GitHubAPIURL
		if len(baseURL) > 0 && baseURL[0] != "" {
			url = baseURL[0]
		}
		singleton = &githubClient{
			client:  &http.Client{},
			baseURL: url,
			headers: map[string]string{
				"Authorization": "token " + os.Getenv(constants.GITHUB_TOKEN),
				"Content-Type":  "application/json",
			},
		}
	})
	return singleton
}

// SearchCode fetches search results using a generic request function.
func (c *githubClient) SearchCode(ctx context.Context, query string) (*structures.Data, error) {
	reqURL := fmt.Sprintf("%s?q=%s", c.baseURL, query)
	fmt.Printf("final url %v\n", reqURL)
	dynamicHeaders := map[string]string{
		// Add any dynamic headers here
	}

	// Merge static and dynamic headers
	headers := c.updateHeaders(dynamicHeaders)

	data := structures.Data{}

	// Use factory to get a GET request handler.
	requestHandler, err := httpClient.NewHTTPRequest(constants.HTTP_METHOD_GET)
	if err != nil {
		return nil, fmt.Errorf("error from NewHTTPRequest: %v", err)
	}

	err = requestHandler.DoRequest(ctx, c.client, reqURL, headers, nil, &data)
	if err != nil {
		return nil, fmt.Errorf("error from DoRequest, %v", err)
	}
	return &data, nil
}

func (c *githubClient) updateHeaders(dynamicHeaders map[string]string) map[string]string {
	for k, v := range c.headers {
		dynamicHeaders[k] = v
	}
	return dynamicHeaders
}
