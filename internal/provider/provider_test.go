// ABOUTME: Tests for the Zabbix Terraform provider configuration.
// ABOUTME: Verifies environment variable fallback and schema validation.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProvider_Configure_EnvironmentVariableFallback(t *testing.T) {
	t.Setenv("ZABBIX_URL", "https://env.example.com/api_jsonrpc.php")
	t.Setenv("ZABBIX_API_TOKEN", "env-token")

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
	t.Setenv("ZABBIX_URL", "")
	t.Setenv("ZABBIX_API_TOKEN", "")

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
	t.Setenv("ZABBIX_URL", "https://env.example.com/api_jsonrpc.php")
	t.Setenv("ZABBIX_API_TOKEN", "env-token")

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
