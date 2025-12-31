// ABOUTME: Terraform data source for looking up existing Zabbix templates.
// ABOUTME: Retrieves template metadata and exported content by technical name.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/terraform-provider-zabbix/internal/zabbix"
)

var _ datasource.DataSource = &TemplateDataSource{}

// TemplateDataSource defines the data source implementation.
type TemplateDataSource struct {
	client *zabbix.Client
}

// TemplateDataSourceModel describes the data source data model.
type TemplateDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Host            types.String `tfsdk:"host"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	UUID            types.String `tfsdk:"uuid"`
	Groups          types.List   `tfsdk:"groups"`
	Tags            types.List   `tfsdk:"tags"`
	ExportedContent types.String `tfsdk:"exported_content"`
}

// NewTemplateDataSource creates a new data source instance.
func NewTemplateDataSource() datasource.DataSource {
	return &TemplateDataSource{}
}

func (d *TemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *TemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Zabbix template by technical name. Returns template metadata and exported content.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the template (templateid in Zabbix).",
				Computed:    true,
			},
			"host": schema.StringAttribute{
				Description: "Technical name of the template to look up.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Visible name of the template.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the template.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "Universally unique identifier of the template.",
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of host group IDs the template belongs to.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tags": schema.ListNestedAttribute{
				Description: "Template tags.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							Description: "Tag name.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value.",
							Computed:    true,
						},
					},
				},
			},
			"exported_content": schema.StringAttribute{
				Description: "Exported template content in YAML format.",
				Computed:    true,
			},
		},
	}
}

func (d *TemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := d.client.GetTemplateByHost(ctx, data.Host.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template",
			fmt.Sprintf("Could not read template with host %q: %s", data.Host.ValueString(), err),
		)
		return
	}

	if template == nil {
		resp.Diagnostics.AddError(
			"Template Not Found",
			fmt.Sprintf("No template found with technical name %q.", data.Host.ValueString()),
		)
		return
	}

	// Export the template content
	exported, err := d.client.ExportConfiguration(ctx, "yaml", []string{template.TemplateID})
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Exporting Template",
			fmt.Sprintf("Could not export template content: %s", err),
		)
	}

	diags := d.apiToModel(ctx, template, &data, exported)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// apiToModel converts the Zabbix API struct to Terraform model.
func (d *TemplateDataSource) apiToModel(ctx context.Context, template *zabbix.Template, data *TemplateDataSourceModel, exportedContent string) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(template.TemplateID)
	data.Host = types.StringValue(template.Host)
	data.Name = types.StringValue(template.Name)
	data.Description = types.StringValue(template.Description)
	data.UUID = types.StringValue(template.UUID)

	// Convert groups
	groupIDs := make([]attr.Value, len(template.Groups))
	for i, g := range template.Groups {
		groupIDs[i] = types.StringValue(g.GroupID)
	}
	groupsList, diagsGroups := types.ListValue(types.StringType, groupIDs)
	diags.Append(diagsGroups...)
	data.Groups = groupsList

	// Convert tags
	tagType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tag":   types.StringType,
			"value": types.StringType,
		},
	}
	if len(template.Tags) > 0 {
		tagValues := make([]attr.Value, len(template.Tags))
		for i, tag := range template.Tags {
			obj, diagsTag := types.ObjectValue(tagType.AttrTypes, map[string]attr.Value{
				"tag":   types.StringValue(tag.Tag),
				"value": types.StringValue(tag.Value),
			})
			diags.Append(diagsTag...)
			tagValues[i] = obj
		}
		tagsList, diagsTags := types.ListValue(tagType, tagValues)
		diags.Append(diagsTags...)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(tagType)
	}

	// Set exported content
	if exportedContent != "" {
		data.ExportedContent = types.StringValue(exportedContent)
	} else {
		data.ExportedContent = types.StringNull()
	}

	return diags
}
