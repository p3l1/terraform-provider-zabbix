// ABOUTME: Zabbix Terraform provider implementation using terraform-plugin-framework.
// ABOUTME: Handles provider configuration (URL and API token) and resource registration.

package provider

import (
	"context"
	"os"

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
				Description: "The URL of the Zabbix API endpoint (e.g., https://zabbix.example.com/api_jsonrpc.php). Can also be set via ZABBIX_URL environment variable.",
				Optional:    true,
			},
			"api_token": schema.StringAttribute{
				Description: "The API token for authenticating with the Zabbix API. Can also be set via ZABBIX_API_TOKEN environment variable.",
				Optional:    true,
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

	url := os.Getenv("ZABBIX_URL")
	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	apiToken := os.Getenv("ZABBIX_API_TOKEN")
	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Missing URL Configuration",
			"The provider requires a URL to be set. "+
				"Set the url attribute in the provider configuration or use the ZABBIX_URL environment variable.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"The provider requires an API token to be set. "+
				"Set the api_token attribute in the provider configuration or use the ZABBIX_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Create Zabbix client with url and apiToken
}

func (p *ZabbixProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ZabbixProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
