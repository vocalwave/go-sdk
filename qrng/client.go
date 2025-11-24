package qrng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the QRNG API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// EntropyResult represents the result of entropy generation
type EntropyResult struct {
	Data          string                 `json:"data"`
	ProofID       string                 `json:"proofId"`
	Signature     string                 `json:"signature"`
	PublicKey     string                 `json:"publicKey"`
	SignatureType string                 `json:"signatureType"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// HealthStatus represents system health
type HealthStatus struct {
	Status    string                 `json:"status"`
	Metrics   map[string]interface{} `json:"metrics"`
	Timestamp string                 `json:"timestamp"`
}

// GenerateOptions options for Generate
type GenerateOptions struct {
	Bytes         int
	Format        string
	Method        string
	SignatureType string
}

// NewClient creates a new QRNG API client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://qrngapi.com",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate generates random entropy
func (c *Client) Generate(opts *GenerateOptions) (*EntropyResult, error) {
	if opts == nil {
		opts = &GenerateOptions{
			Bytes:  32,
			Format: "hex",
		}
	}

	params := url.Values{}
	params.Add("bytes", fmt.Sprintf("%d", opts.Bytes))
	params.Add("format", opts.Format)
	if opts.Method != "" {
		params.Add("method", opts.Method)
	}
	if opts.SignatureType != "" {
		params.Add("signatureType", opts.SignatureType)
	}

	reqURL := fmt.Sprintf("%s/api/random?%s", c.BaseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.Unmarshal(body, &errResp)
		if msg, ok := errResp["error"].(string); ok {
			return nil, fmt.Errorf("API error: %s", msg)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result EntropyResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// Health gets system health status
func (c *Client) Health() (*HealthStatus, error) {
	reqURL := fmt.Sprintf("%s/api/health", c.BaseURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var health HealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &health, nil
}
