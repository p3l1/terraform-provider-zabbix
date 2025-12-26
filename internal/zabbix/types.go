// ABOUTME: Defines JSON-RPC 2.0 request/response types and error handling for Zabbix API.
// ABOUTME: Contains shared types used across the Zabbix client package.

package zabbix

import (
	"encoding/json"
	"fmt"
)

// Request represents a JSON-RPC 2.0 request to the Zabbix API.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
	Auth    string      `json:"auth,omitempty"`
}

// Response represents a JSON-RPC 2.0 response from the Zabbix API.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
	ID      int             `json:"id"`
}

// Error represents a Zabbix API error response.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (e *Error) Error() string {
	if e.Data != "" {
		return fmt.Sprintf("zabbix api error %d: %s - %s", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("zabbix api error %d: %s", e.Code, e.Message)
}

// APIError wraps a Zabbix API error with additional context.
type APIError struct {
	Method string
	Err    *Error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("method %s: %s", e.Method, e.Err.Error())
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// HTTPError represents an HTTP-level error when communicating with the Zabbix API.
type HTTPError struct {
	StatusCode int
	Status     string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("zabbix api http error: %s", e.Status)
}
