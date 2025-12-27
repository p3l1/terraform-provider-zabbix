// ABOUTME: Unit tests for host API methods using mock HTTP responses.
// ABOUTME: Tests cover CRUD operations and error handling for hosts.

package zabbix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "host.create" {
			t.Errorf("expected method 'host.create', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["host"] != "test-server" {
			t.Errorf("expected host 'test-server', got '%v'", params["host"])
		}
		if params["name"] != "Test Server" {
			t.Errorf("expected name 'Test Server', got '%v'", params["name"])
		}

		groups, ok := params["groups"].([]interface{})
		if !ok || len(groups) != 1 {
			t.Fatalf("expected groups to be array with 1 element, got %v", params["groups"])
		}
		group := groups[0].(map[string]interface{})
		if group["groupid"] != "2" {
			t.Errorf("expected groupid '2', got '%v'", group["groupid"])
		}

		interfaces, ok := params["interfaces"].([]interface{})
		if !ok || len(interfaces) != 1 {
			t.Fatalf("expected interfaces to be array with 1 element, got %v", params["interfaces"])
		}
		iface := interfaces[0].(map[string]interface{})
		if iface["type"] != float64(1) {
			t.Errorf("expected interface type 1, got '%v'", iface["type"])
		}
		if iface["ip"] != "192.168.1.100" {
			t.Errorf("expected ip '192.168.1.100', got '%v'", iface["ip"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": ["10084"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		Host:   "test-server",
		Name:   "Test Server",
		Status: 0,
		Groups: []HostGroupID{{GroupID: "2"}},
		Interfaces: []HostInterface{{
			Type:  1,
			Main:  1,
			UseIP: 1,
			IP:    "192.168.1.100",
			DNS:   "",
			Port:  "10050",
		}},
	}
	hostID, err := client.CreateHost(context.Background(), host)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hostID != "10084" {
		t.Errorf("expected hostID '10084', got '%s'", hostID)
	}
}

func TestCreateHost_WithTemplates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		templates, ok := params["templates"].([]interface{})
		if !ok || len(templates) != 2 {
			t.Fatalf("expected templates to be array with 2 elements, got %v", params["templates"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": ["10084"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		Host:   "test-server",
		Groups: []HostGroupID{{GroupID: "2"}},
		Templates: []TemplateID{
			{TemplateID: "10001"},
			{TemplateID: "10002"},
		},
		Interfaces: []HostInterface{{
			Type:  1,
			Main:  1,
			UseIP: 1,
			IP:    "192.168.1.100",
			Port:  "10050",
		}},
	}
	hostID, err := client.CreateHost(context.Background(), host)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hostID != "10084" {
		t.Errorf("expected hostID '10084', got '%s'", hostID)
	}
}

func TestCreateHost_WithTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		tags, ok := params["tags"].([]interface{})
		if !ok || len(tags) != 1 {
			t.Fatalf("expected tags to be array with 1 element, got %v", params["tags"])
		}
		tag := tags[0].(map[string]interface{})
		if tag["tag"] != "environment" {
			t.Errorf("expected tag 'environment', got '%v'", tag["tag"])
		}
		if tag["value"] != "production" {
			t.Errorf("expected value 'production', got '%v'", tag["value"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": ["10084"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		Host:   "test-server",
		Groups: []HostGroupID{{GroupID: "2"}},
		Tags: []HostTag{{
			Tag:   "environment",
			Value: "production",
		}},
		Interfaces: []HostInterface{{
			Type:  1,
			Main:  1,
			UseIP: 1,
			IP:    "192.168.1.100",
			Port:  "10050",
		}},
	}
	hostID, err := client.CreateHost(context.Background(), host)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hostID != "10084" {
		t.Errorf("expected hostID '10084', got '%s'", hostID)
	}
}

func TestCreateHost_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		Host:   "test-server",
		Groups: []HostGroupID{{GroupID: "2"}},
		Interfaces: []HostInterface{{
			Type:  1,
			Main:  1,
			UseIP: 1,
			IP:    "192.168.1.100",
			Port:  "10050",
		}},
	}
	_, err := client.CreateHost(context.Background(), host)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateHost_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "Host already exists.",
			},
			ID: 1,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		Host:   "test-server",
		Groups: []HostGroupID{{GroupID: "2"}},
		Interfaces: []HostInterface{{
			Type:  1,
			Main:  1,
			UseIP: 1,
			IP:    "192.168.1.100",
			Port:  "10050",
		}},
	}
	_, err := client.CreateHost(context.Background(), host)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "host.create" {
		t.Errorf("expected method 'host.create', got '%s'", apiErr.Method)
	}
}

func TestGetHost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "host.get" {
			t.Errorf("expected method 'host.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		hostIDs, ok := params["hostids"].([]interface{})
		if !ok || len(hostIDs) != 1 || hostIDs[0] != "10084" {
			t.Errorf("expected hostids ['10084'], got '%v'", params["hostids"])
		}

		if params["selectGroups"] != "extend" {
			t.Errorf("expected selectGroups 'extend', got '%v'", params["selectGroups"])
		}
		if params["selectInterfaces"] != "extend" {
			t.Errorf("expected selectInterfaces 'extend', got '%v'", params["selectInterfaces"])
		}
		if params["selectTags"] != "extend" {
			t.Errorf("expected selectTags 'extend', got '%v'", params["selectTags"])
		}
		if params["selectParentTemplates"] != "extend" {
			t.Errorf("expected selectParentTemplates 'extend', got '%v'", params["selectParentTemplates"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result: json.RawMessage(`[{
				"hostid": "10084",
				"host": "test-server",
				"name": "Test Server",
				"status": "0",
				"groups": [{"groupid": "2", "name": "Linux servers"}],
				"interfaces": [{"interfaceid": "1", "type": "1", "main": "1", "useip": "1", "ip": "192.168.1.100", "dns": "", "port": "10050"}],
				"tags": [{"tag": "environment", "value": "production"}],
				"parentTemplates": [{"templateid": "10001", "name": "Template OS Linux"}]
			}]`),
			ID: req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host, err := client.GetHost(context.Background(), "10084")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host == nil {
		t.Fatal("expected host, got nil")
	}
	if host.HostID != "10084" {
		t.Errorf("expected hostid '10084', got '%s'", host.HostID)
	}
	if host.Host != "test-server" {
		t.Errorf("expected host 'test-server', got '%s'", host.Host)
	}
	if host.Name != "Test Server" {
		t.Errorf("expected name 'Test Server', got '%s'", host.Name)
	}
	if host.Status != 0 {
		t.Errorf("expected status 0, got %d", host.Status)
	}
	if len(host.Groups) != 1 || host.Groups[0].GroupID != "2" {
		t.Errorf("expected groups with groupid '2', got %v", host.Groups)
	}
	if len(host.Interfaces) != 1 || host.Interfaces[0].IP != "192.168.1.100" {
		t.Errorf("expected interface with IP '192.168.1.100', got %v", host.Interfaces)
	}
	if len(host.Tags) != 1 || host.Tags[0].Tag != "environment" {
		t.Errorf("expected tag 'environment', got %v", host.Tags)
	}
	if len(host.ParentTemplates) != 1 || host.ParentTemplates[0].TemplateID != "10001" {
		t.Errorf("expected template with id '10001', got %v", host.ParentTemplates)
	}
}

func TestGetHost_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host, err := client.GetHost(context.Background(), "99999")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != nil {
		t.Errorf("expected nil host, got %v", host)
	}
}

func TestGetHostByName_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "host.get" {
			t.Errorf("expected method 'host.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		filter, ok := params["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected filter to be a map, got %T", params["filter"])
		}
		if filter["host"] != "test-server" {
			t.Errorf("expected filter host 'test-server', got '%v'", filter["host"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result: json.RawMessage(`[{
				"hostid": "10084",
				"host": "test-server",
				"name": "Test Server",
				"status": "0",
				"groups": [{"groupid": "2", "name": "Linux servers"}],
				"interfaces": [{"interfaceid": "1", "type": "1", "main": "1", "useip": "1", "ip": "192.168.1.100", "dns": "", "port": "10050"}],
				"tags": [],
				"parentTemplates": []
			}]`),
			ID: req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host, err := client.GetHostByName(context.Background(), "test-server")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host == nil {
		t.Fatal("expected host, got nil")
	}
	if host.HostID != "10084" {
		t.Errorf("expected hostid '10084', got '%s'", host.HostID)
	}
	if host.Host != "test-server" {
		t.Errorf("expected host 'test-server', got '%s'", host.Host)
	}
}

func TestGetHostByName_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host, err := client.GetHostByName(context.Background(), "nonexistent")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != nil {
		t.Errorf("expected nil host, got %v", host)
	}
}

func TestUpdateHost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "host.update" {
			t.Errorf("expected method 'host.update', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["hostid"] != "10084" {
			t.Errorf("expected hostid '10084', got '%v'", params["hostid"])
		}
		if params["name"] != "Updated Server" {
			t.Errorf("expected name 'Updated Server', got '%v'", params["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": ["10084"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		HostID: "10084",
		Name:   "Updated Server",
	}
	err := client.UpdateHost(context.Background(), host)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateHost_WithGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		groups, ok := params["groups"].([]interface{})
		if !ok || len(groups) != 2 {
			t.Fatalf("expected groups to be array with 2 elements, got %v", params["groups"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": ["10084"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		HostID: "10084",
		Groups: []HostGroupID{
			{GroupID: "2"},
			{GroupID: "5"},
		},
	}
	err := client.UpdateHost(context.Background(), host)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateHost_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	host := &Host{
		HostID: "10084",
		Name:   "Updated Server",
	}
	err := client.UpdateHost(context.Background(), host)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteHost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "host.delete" {
			t.Errorf("expected method 'host.delete', got '%s'", req.Method)
		}

		params, ok := req.Params.([]interface{})
		if !ok {
			t.Fatalf("expected params to be an array, got %T", req.Params)
		}
		if len(params) != 1 || params[0] != "10084" {
			t.Errorf("expected params ['10084'], got '%v'", params)
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": ["10084"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteHost(context.Background(), "10084")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteHost_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"hostids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteHost(context.Background(), "10084")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteHost_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "No permissions to referred object or it does not exist!",
			},
			ID: 1,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteHost(context.Background(), "99999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "host.delete" {
		t.Errorf("expected method 'host.delete', got '%s'", apiErr.Method)
	}
}
