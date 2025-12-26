// ABOUTME: Integration tests for the Zabbix API client against a real Zabbix instance.
// ABOUTME: Requires TF_ACC=1 and a running Docker Zabbix environment.

package zabbix

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	defaultTestURL   = "http://127.0.0.1:8080/api_jsonrpc.php"
	defaultTestToken = "071fb9d2e8f72cf9c40128f0f5aab3def1bab0893413314b083fdcb4551eb01a"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run integration tests")
	}

	url := os.Getenv("ZABBIX_URL")
	if url == "" {
		url = defaultTestURL
	}

	token := os.Getenv("ZABBIX_API_TOKEN")
	if token == "" {
		token = defaultTestToken
	}

	return NewClient(url, token)
}

func TestIntegration_APIVersion(t *testing.T) {
	client := newTestClient(t)
	result, err := client.Request("apiinfo.version", nil)

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var version string
	if err := json.Unmarshal(result, &version); err != nil {
		t.Fatalf("Failed to unmarshal version: %v", err)
	}

	if version == "" {
		t.Error("Expected non-empty version string")
	}

	t.Logf("Zabbix API version: %s", version)
}

func TestIntegration_HostGet(t *testing.T) {
	client := newTestClient(t)
	result, err := client.Request("host.get", map[string]interface{}{
		"output": []string{"hostid", "host"},
	})

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var hosts []map[string]interface{}
	if err := json.Unmarshal(result, &hosts); err != nil {
		t.Fatalf("Failed to unmarshal hosts: %v", err)
	}

	t.Logf("Found %d hosts", len(hosts))
}

func TestIntegration_InvalidToken(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run integration tests")
	}

	url := os.Getenv("ZABBIX_URL")
	if url == "" {
		url = defaultTestURL
	}

	client := NewClient(url, "invalid-token")
	_, err := client.Request("host.get", nil)

	if err == nil {
		t.Fatal("Expected error with invalid token")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T: %v", err, err)
	}

	t.Logf("Got expected error: %v", apiErr)
}
