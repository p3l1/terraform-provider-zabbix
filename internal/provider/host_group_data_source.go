// ABOUTME: Terraform data source for looking up existing Zabbix host groups.
// ABOUTME: Retrieves host group information by name.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/zabbix-terraform/internal/zabbix"
)

var _ datasource.DataSource = &HostGroupDataSource{}

// HostGroupDataSource defines the data source implementation.
type HostGroupDataSource struct {
	client *zabbix.Client
}

// HostGroupDataSourceModel describes the data source data model.
type HostGroupDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

// NewHostGroupDataSource creates a new data source instance.
func NewHostGroupDataSource() datasource.DataSource {
	return &HostGroupDataSource{}
}

func (d *HostGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_group"
}

func (d *HostGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Zabbix host group by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the host group (groupid in Zabbix).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the host group to look up.",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The universally unique identifier of the host group.",
				Computed:    true,
			},
		},
	}
}

func (d *HostGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetHostGroupByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host Group",
			fmt.Sprintf("Could not read host group with name %q: %s", data.Name.ValueString(), err),
		)
		return
	}

	if group == nil {
		resp.Diagnostics.AddError(
			"Host Group Not Found",
			fmt.Sprintf("No host group found with name %q.", data.Name.ValueString()),
		)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
