package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	httpClient "github.com/Thesohan/weaveGitHubSearchService/server/http_client"
)

const (
	defaultGitHubAPIURL = "https://api.github.com/search/code"
	githubTokenEnvName  = "GITHUB_TOKEN"
	httpClientTimeout   = 3 * time.Second
)

var (
	once      sync.Once
	singleton CodeSearcher
)

type searchResponse struct {
	Items []struct {
		HTMLURL    string `json:"html_url"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	} `json:"items"`
}

type githubClient struct {
	client            *http.Client
	baseURL           string
	headers           map[string]string
	maximumRetryDelay *time.Duration
}

// IGithubClient interface for API requests.
type CodeSearcher interface {
	SearchCode(ctx context.Context, query string) (*searchResponse, error)
}

type GitHubClientOption func(*githubClient)

func WithBaseURL(baseURL string) GitHubClientOption {
	return func(c *githubClient) { c.baseURL = baseURL }
}

func WithToken(token string) GitHubClientOption {
	return func(c *githubClient) { c.headers["Authorization"] = "token " + token }
}

func WithMaximumRetryDelay(maximumRetryDelay *time.Duration) GitHubClientOption {
	return func(c *githubClient) { c.maximumRetryDelay = maximumRetryDelay }
}

func NewGitHubClient(opts ...GitHubClientOption) CodeSearcher {
	client := &githubClient{
		client:  &http.Client{Timeout: httpClientTimeout},
		baseURL: defaultGitHubAPIURL,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for _, opt := range opts {
		opt(client)
	}
	if _, ok := client.headers["Authorization"]; !ok {
		client.headers["Authorization"] = "token " + os.Getenv(githubTokenEnvName)
	}
	return client
}

// SearchCode fetches search results using a generic request function.
func (c *githubClient) SearchCode(ctx context.Context, query string) (*searchResponse, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	param := url.Values{}
	param.Add("q", query)
	base.RawQuery = param.Encode()

	reqURL := base.String()
	// Use factory to get a GET request handler.
	data := searchResponse{}
	err = httpClient.DoRequest(ctx, c.client, http.MethodGet, reqURL, c.headers, nil /* body */, &data, c.maximumRetryDelay)
	if err != nil {
		return nil, fmt.Errorf("error from DoRequest, %v", err)
	}
	return &data, nil
}
