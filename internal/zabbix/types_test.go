// ABOUTME: Unit tests for Zabbix API types and error formatting.
// ABOUTME: Tests cover error message formatting and error wrapping behavior.

package zabbix

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "error with data",
			err: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "No permissions",
			},
			expected: "zabbix api error -32602: Invalid params. - No permissions",
		},
		{
			name: "error without data",
			err: &Error{
				Code:    -32600,
				Message: "Invalid request.",
			},
			expected: "zabbix api error -32600: Invalid request.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	apiErr := &APIError{
		Method: "host.get",
		Err: &Error{
			Code:    -32602,
			Message: "Invalid params.",
			Data:    "No permissions",
		},
	}

	expected := "method host.get: zabbix api error -32602: Invalid params. - No permissions"
	if got := apiErr.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	innerErr := &Error{
		Code:    -32602,
		Message: "Invalid params.",
	}
	apiErr := &APIError{
		Method: "host.get",
		Err:    innerErr,
	}

	unwrapped := apiErr.Unwrap()
	if unwrapped != innerErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}

	if !errors.Is(apiErr, innerErr) {
		t.Error("errors.Is should return true for wrapped error")
	}
}

func TestHTTPError_Error(t *testing.T) {
	httpErr := &HTTPError{
		StatusCode: 500,
		Status:     "500 Internal Server Error",
	}

	expected := "zabbix api http error: 500 Internal Server Error"
	if got := httpErr.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}
