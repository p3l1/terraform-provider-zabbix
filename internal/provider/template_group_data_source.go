// ABOUTME: Terraform data source for looking up existing Zabbix template groups.
// ABOUTME: Retrieves template group information by name.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/terraform-provider-zabbix/internal/zabbix"
)

var _ datasource.DataSource = &TemplateGroupDataSource{}

// TemplateGroupDataSource defines the data source implementation.
type TemplateGroupDataSource struct {
	client *zabbix.Client
}

// TemplateGroupDataSourceModel describes the data source data model.
type TemplateGroupDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

// NewTemplateGroupDataSource creates a new data source instance.
func NewTemplateGroupDataSource() datasource.DataSource {
	return &TemplateGroupDataSource{}
}

func (d *TemplateGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_group"
}

func (d *TemplateGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Zabbix template group by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the template group (groupid in Zabbix).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the template group to look up.",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The universally unique identifier of the template group.",
				Computed:    true,
			},
		},
	}
}

func (d *TemplateGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*zabbix.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *zabbix.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *TemplateGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplateGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetTemplateGroupByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template Group",
			fmt.Sprintf("Could not read template group with name %q: %s", data.Name.ValueString(), err),
		)
		return
	}

	if group == nil {
		resp.Diagnostics.AddError(
			"Template Group Not Found",
			fmt.Sprintf("No template group found with name %q.", data.Name.ValueString()),
		)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
