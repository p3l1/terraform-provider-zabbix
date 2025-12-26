// ABOUTME: HTTP client for communicating with the Zabbix JSON-RPC 2.0 API.
// ABOUTME: Handles authentication, request serialization, and response parsing.

package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Client is a Zabbix API client.
type Client struct {
	URL        string
	Token      string
	HTTPClient *http.Client
	requestID  atomic.Int64
}

// NewClient creates a new Zabbix API client with default settings.
func NewClient(url, token string) *Client {
	return &Client{
		URL:   url,
		Token: token,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithTimeout creates a new Zabbix API client with a custom timeout.
func NewClientWithTimeout(url, token string, timeout time.Duration) *Client {
	return &Client{
		URL:   url,
		Token: token,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Methods that don't require authentication.
var noAuthMethods = map[string]bool{
	"apiinfo.version": true,
}

// Request sends a JSON-RPC 2.0 request to the Zabbix API.
func (c *Client) Request(method string, params interface{}) (json.RawMessage, error) {
	if params == nil {
		params = map[string]interface{}{}
	}

	req := Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      int(c.requestID.Add(1)),
	}

	if !noAuthMethods[method] {
		req.Auth = c.Token
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json-rpc")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			StatusCode: httpResp.StatusCode,
			Status:     httpResp.Status,
		}
	}

	var resp Response
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.Error != nil {
		return nil, &APIError{
			Method: method,
			Err:    resp.Error,
		}
	}

	return resp.Result, nil
}
