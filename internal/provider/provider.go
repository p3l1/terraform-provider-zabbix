// ABOUTME: Zabbix Terraform provider implementation using terraform-plugin-framework.
// ABOUTME: Handles provider configuration (URL and API token) and resource registration.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &ZabbixProvider{}

// ZabbixProvider implements the Zabbix Terraform provider.
type ZabbixProvider struct {
	version string
}

// ZabbixProviderModel describes the provider configuration data.
type ZabbixProviderModel struct {
	URL      types.String `tfsdk:"url"`
	APIToken types.String `tfsdk:"api_token"`
}

// New creates a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZabbixProvider{
			version: version,
		}
	}
}

func (p *ZabbixProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zabbix"
	resp.Version = p.version
}

func (p *ZabbixProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing Zabbix monitoring infrastructure.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The URL of the Zabbix API endpoint (e.g., https://zabbix.example.com/api_jsonrpc.php).",
				Required:    true,
			},
			"api_token": schema.StringAttribute{
				Description: "The API token for authenticating with the Zabbix API.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *ZabbixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ZabbixProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration validation will be added when we implement the Zabbix client
}

func (p *ZabbixProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ZabbixProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
