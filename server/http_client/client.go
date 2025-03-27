package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const defaultMaximumRetryDelay = 30 * time.Second

// DoRequest performs an HTTP request with retries, request body handling, and response decoding.
func DoRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body interface{}, response interface{}, maximumRetryDelay *time.Duration) error {
	var requestBody io.Reader

	// Encode request body if provided
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewReader(data)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return executeRequest(client, req, response, maximumRetryDelay)
}

// executeRequest performs the HTTP request with retries and handles response parsing.
func executeRequest(client *http.Client, req *http.Request, response interface{}, maximumRetryDelay *time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = defaultMaximumRetryDelay
	if maximumRetryDelay != nil {
		bo.MaxElapsedTime = *maximumRetryDelay
	}

	operation := func() error {
		log.Printf("HTTP %s request to %s at: %v", req.Method, req.URL, time.Now().Format(time.RFC3339))

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("API error: %s", resp.Status)
		}

		if response != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			if err := json.Unmarshal(bodyBytes, response); err != nil {
				return fmt.Errorf("failed to unmarshal response body: %w", err)
			}
			return nil
		}
		return nil
	}

	if err := backoff.Retry(operation, bo); err != nil {
		return fmt.Errorf("request failed after retries: %w", err)
	}
	return nil
}
