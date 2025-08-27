package am

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

type APIClient struct {
	Core
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(name, baseURL string, opts ...Option) *APIClient {
	core := NewCore(name, opts...)
	client := &APIClient{
		Core:       core,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}

	return client
}

func (c *APIClient) request(method, path string, body, target interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("errir creating API request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	c.Log().Debugf("API Request: %s %s", method, url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.Log().Errorf("API Request failed: %v", err)
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	c.Log().Debugf("API Response: %d %s", resp.StatusCode, string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("api request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResponse Response
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return fmt.Errorf("error decoding API response structure: %w", err)
	}

	if apiResponse.Status != StatusSuccess {
		return fmt.Errorf("api returned error: %s", apiResponse.Message)
	}

	if target != nil {
		dataBytes, err := json.Marshal(apiResponse.Data)
		if err != nil {
			return fmt.Errorf("error marshaling API response data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, target); err != nil {
			return fmt.Errorf("error unmarshaling data into target: %w", err)
		}
	}

	return nil
}

// Get sends a GET request to the specified path.
func (c *APIClient) Get(path string, target interface{}) error {
	return c.request(http.MethodGet, path, nil, target)
}

// Post sends a POST request to the specified path with the given body.
func (c *APIClient) Post(path string, body, target interface{}) error {
	return c.request(http.MethodPost, path, body, target)
}

// Put sends a PUT request to the specified path with the given body.
func (c *APIClient) Put(path string, body, target interface{}) error {
	return c.request(http.MethodPut, path, body, target)
}

// Delete sends a DELETE request to the specified path.
func (c *APIClient) Delete(path string) error {
	return c.request(http.MethodDelete, path, nil, nil)
}
