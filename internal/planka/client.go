package planka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Planka API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// LoginResponse represents the response from a login request
type LoginResponse struct {
	Item string `json:"item"` // The access token
}

// NewClient creates a new Planka API client with a token
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithPassword creates a new Planka API client by logging in with username/password
func NewClientWithPassword(baseURL, username, password string) (*Client, error) {
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	loginReq := map[string]string{
		"emailOrUsername": username,
		"password":         password,
	}

	var loginResp LoginResponse
	if err := client.postWithoutAuth("/api/access-tokens", loginReq, &loginResp); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	client.token = loginResp.Item
	return client, nil
}

// postWithoutAuth performs a POST request without authentication (for login)
func (c *Client) postWithoutAuth(endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// doRequest performs an HTTP request to the Planka API
func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// get performs a GET request
func (c *Client) get(endpoint string, result interface{}) error {
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		// Read the body first to check if it's valid JSON
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		
		// Check if response is HTML (starts with <)
		if len(bodyBytes) > 0 && bodyBytes[0] == '<' {
			return fmt.Errorf("received HTML instead of JSON for endpoint %s. Response preview: %s", endpoint, string(bodyBytes[:min(200, len(bodyBytes))]))
		}
		
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("failed to decode JSON response for endpoint %s: %w. Response preview: %s", endpoint, err, string(bodyBytes[:min(200, len(bodyBytes))]))
		}
	}

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// post performs a POST request
func (c *Client) post(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.doRequest("POST", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// patch performs a PATCH request
func (c *Client) patch(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.doRequest("PATCH", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// delete performs a DELETE request
func (c *Client) delete(endpoint string) error {
	resp, err := c.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

