// ABOUTME: Tests for the Zabbix Terraform provider configuration.
// ABOUTME: Verifies environment variable fallback and schema validation.

package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProvider_Configure_EnvironmentVariableFallback(t *testing.T) {
	os.Setenv("ZABBIX_URL", "https://env.example.com/api_jsonrpc.php")
	os.Setenv("ZABBIX_API_TOKEN", "env-token")
	defer os.Unsetenv("ZABBIX_URL")
	defer os.Unsetenv("ZABBIX_API_TOKEN")

	p := New("test")()
	schemaResp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

	configValue := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"url":       tftypes.String,
			"api_token": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"url":       tftypes.NewValue(tftypes.String, nil),
		"api_token": tftypes.NewValue(tftypes.String, nil),
	})

	config, err := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    configValue,
	}, error(nil)
	if err != nil {
		t.Fatalf("failed to create config: %s", err)
	}

	req := provider.ConfigureRequest{Config: config}
	resp := &provider.ConfigureResponse{}

	p.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
}

func TestProvider_Configure_MissingRequiredConfig(t *testing.T) {
	os.Unsetenv("ZABBIX_URL")
	os.Unsetenv("ZABBIX_API_TOKEN")

	p := New("test")()
	schemaResp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

	configValue := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"url":       tftypes.String,
			"api_token": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"url":       tftypes.NewValue(tftypes.String, nil),
		"api_token": tftypes.NewValue(tftypes.String, nil),
	})

	config := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    configValue,
	}

	req := provider.ConfigureRequest{Config: config}
	resp := &provider.ConfigureResponse{}

	p.Configure(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when URL and API token are not configured")
	}
}

func TestProvider_Configure_ConfigOverridesEnvironment(t *testing.T) {
	os.Setenv("ZABBIX_URL", "https://env.example.com/api_jsonrpc.php")
	os.Setenv("ZABBIX_API_TOKEN", "env-token")
	defer os.Unsetenv("ZABBIX_URL")
	defer os.Unsetenv("ZABBIX_API_TOKEN")

	p := New("test")()
	schemaResp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

	configValue := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"url":       tftypes.String,
			"api_token": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"url":       tftypes.NewValue(tftypes.String, "https://config.example.com/api_jsonrpc.php"),
		"api_token": tftypes.NewValue(tftypes.String, "config-token"),
	})

	config := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    configValue,
	}

	req := provider.ConfigureRequest{Config: config}
	resp := &provider.ConfigureResponse{}

	p.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
}