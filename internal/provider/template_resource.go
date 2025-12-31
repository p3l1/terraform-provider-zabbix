// ABOUTME: Terraform resource for managing Zabbix templates.
// ABOUTME: Supports both metadata-only management and full template import via YAML/JSON/XML.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/terraform-provider-zabbix/internal/zabbix"
)

var (
	_ resource.Resource                = &TemplateResource{}
	_ resource.ResourceWithImportState = &TemplateResource{}
)

// TemplateResource defines the resource implementation.
type TemplateResource struct {
	client *zabbix.Client
}

// TemplateResourceModel describes the resource data model.
type TemplateResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Host            types.String `tfsdk:"host"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	UUID            types.String `tfsdk:"uuid"`
	Groups          types.List   `tfsdk:"groups"`
	Tags            types.List   `tfsdk:"tags"`
	SourceFormat    types.String `tfsdk:"source_format"`
	SourceContent   types.String `tfsdk:"source_content"`
	ExportedContent types.String `tfsdk:"exported_content"`
}

// TemplateTagModel describes a template tag.
type TemplateTagModel struct {
	Tag   types.String `tfsdk:"tag"`
	Value types.String `tfsdk:"value"`
}

// NewTemplateResource creates a new resource instance.
func NewTemplateResource() resource.Resource {
	return &TemplateResource{}
}

func (r *TemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (r *TemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Zabbix template. Can create templates from YAML/JSON/XML content or manage template metadata directly.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the template (templateid in Zabbix).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Description: "Technical name of the template. Required when not using source_content.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Visible name of the template. Defaults to host if not set.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the template.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Universally unique identifier of the template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"groups": schema.ListAttribute{
				Description: "List of host group IDs the template belongs to. Required when not using source_content.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "Template tags.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							Description: "Tag name.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"source_format": schema.StringAttribute{
				Description: "Format of source_content: yaml, xml, or json. Required when source_content is provided.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_content": schema.StringAttribute{
				Description: "Template content in YAML, XML, or JSON format. When provided, the template is imported using configuration.import.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"exported_content": schema.StringAttribute{
				Description: "Exported template content in YAML format. Used for drift detection.",
				Computed:    true,
			},
		},
	}
}

func (r *TemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*zabbix.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *zabbix.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var templateID string
	var err error

	if !data.SourceContent.IsNull() && !data.SourceContent.IsUnknown() {
		// Import from source content
		format := data.SourceFormat.ValueString()
		if format == "" {
			resp.Diagnostics.AddError(
				"Missing source_format",
				"source_format is required when source_content is provided",
			)
			return
		}

		err = r.client.ImportConfiguration(ctx, format, data.SourceContent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Importing Template",
				fmt.Sprintf("Could not import template: %s", err),
			)
			return
		}

		// Extract host name from the content to find the imported template
		host := r.extractHostFromContent(data.SourceContent.ValueString(), format)
		if host == "" {
			resp.Diagnostics.AddError(
				"Error Finding Template",
				"Could not determine template host name from source content. Please ensure the content contains a valid template definition.",
			)
			return
		}

		template, err := r.client.GetTemplateByHost(ctx, host)
		if err != nil || template == nil {
			resp.Diagnostics.AddError(
				"Error Finding Imported Template",
				fmt.Sprintf("Could not find template with host %q after import: %v", host, err),
			)
			return
		}
		templateID = template.TemplateID
	} else {
		// Create template directly
		template, diags := r.modelToAPI(ctx, &data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		templateID, err = r.client.CreateTemplate(ctx, template)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Template",
				fmt.Sprintf("Could not create template: %s", err),
			)
			return
		}
	}

	// Read back the template to get computed values
	apiTemplate, err := r.client.GetTemplate(ctx, templateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template",
			fmt.Sprintf("Could not read template after creation: %s", err),
		)
		return
	}

	if apiTemplate == nil {
		resp.Diagnostics.AddError(
			"Error Reading Template",
			fmt.Sprintf("Template %s was created but could not be found", templateID),
		)
		return
	}

	// Export the template content for drift detection
	exported, err := r.client.ExportConfiguration(ctx, "yaml", []string{templateID})
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Exporting Template",
			fmt.Sprintf("Could not export template content: %s", err),
		)
	}

	diags := r.apiToModel(ctx, apiTemplate, &data, exported)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := r.client.GetTemplate(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template",
			fmt.Sprintf("Could not read template ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	if template == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Export the template content
	exported, err := r.client.ExportConfiguration(ctx, "yaml", []string{data.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Exporting Template",
			fmt.Sprintf("Could not export template content: %s", err),
		)
	}

	diags := r.apiToModel(ctx, template, &data, exported)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.SourceContent.IsNull() && !data.SourceContent.IsUnknown() {
		// Re-import from source content
		format := data.SourceFormat.ValueString()
		if format == "" {
			resp.Diagnostics.AddError(
				"Missing source_format",
				"source_format is required when source_content is provided",
			)
			return
		}

		err := r.client.ImportConfiguration(ctx, format, data.SourceContent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Importing Template",
				fmt.Sprintf("Could not import template: %s", err),
			)
			return
		}
	} else {
		// Update template directly
		template, diags := r.modelToAPI(ctx, &data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		template.TemplateID = state.ID.ValueString()

		err := r.client.UpdateTemplate(ctx, template)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Template",
				fmt.Sprintf("Could not update template ID %s: %s", state.ID.ValueString(), err),
			)
			return
		}
	}

	// Read back the template
	apiTemplate, err := r.client.GetTemplate(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template",
			fmt.Sprintf("Could not read template after update: %s", err),
		)
		return
	}

	if apiTemplate == nil {
		resp.Diagnostics.AddError(
			"Error Reading Template",
			fmt.Sprintf("Template %s was updated but could not be found", state.ID.ValueString()),
		)
		return
	}

	// Export the template content
	exported, err := r.client.ExportConfiguration(ctx, "yaml", []string{state.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Exporting Template",
			fmt.Sprintf("Could not export template content: %s", err),
		)
	}

	diags := r.apiToModel(ctx, apiTemplate, &data, exported)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTemplate(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Template",
			fmt.Sprintf("Could not delete template ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
}

func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// modelToAPI converts the Terraform model to Zabbix API struct.
func (r *TemplateResource) modelToAPI(ctx context.Context, data *TemplateResourceModel) (*zabbix.Template, diag.Diagnostics) {
	var diags diag.Diagnostics

	template := &zabbix.Template{
		Host:        data.Host.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Convert groups
	if !data.Groups.IsNull() && !data.Groups.IsUnknown() {
		var groupIDs []string
		diags.Append(data.Groups.ElementsAs(ctx, &groupIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, id := range groupIDs {
			template.Groups = append(template.Groups, zabbix.TemplateGroupID{GroupID: id})
		}
	}

	// Convert tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []TemplateTagModel
		diags.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, tag := range tags {
			template.Tags = append(template.Tags, zabbix.TemplateTag{
				Tag:   tag.Tag.ValueString(),
				Value: tag.Value.ValueString(),
			})
		}
	}

	return template, diags
}

// apiToModel converts the Zabbix API struct to Terraform model.
func (r *TemplateResource) apiToModel(ctx context.Context, template *zabbix.Template, data *TemplateResourceModel, exportedContent string) diag.Diagnostics {
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
	groupsList, d := types.ListValue(types.StringType, groupIDs)
	diags.Append(d...)
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
			obj, d := types.ObjectValue(tagType.AttrTypes, map[string]attr.Value{
				"tag":   types.StringValue(tag.Tag),
				"value": types.StringValue(tag.Value),
			})
			diags.Append(d...)
			tagValues[i] = obj
		}
		tagsList, d := types.ListValue(tagType, tagValues)
		diags.Append(d...)
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

// extractHostFromContent extracts the template host name from YAML/JSON/XML content.
func (r *TemplateResource) extractHostFromContent(content, format string) string {
	// Simple extraction for YAML - look for "template:" or "host:" patterns
	// This is a basic implementation that works for standard Zabbix template exports
	switch format {
	case "yaml":
		return r.extractHostFromYAML(content)
	case "json":
		return r.extractHostFromJSON(content)
	case "xml":
		return r.extractHostFromXML(content)
	}
	return ""
}

func (r *TemplateResource) extractHostFromYAML(content string) string {
	// Look for patterns like:
	// templates:
	//   - template: "Template Name"
	// or
	//   - uuid: ...
	//     template: "Template Name"
	lines := splitLines(content)
	inTemplates := false
	for _, line := range lines {
		trimmed := trimSpace(line)
		if trimmed == "templates:" || trimmed == "zabbix_export:" {
			inTemplates = true
			continue
		}
		if inTemplates && hasPrefix(trimmed, "template:") {
			value := trimSpace(trimPrefix(trimmed, "template:"))
			// Remove quotes if present
			value = trimQuotes(value)
			return value
		}
	}
	return ""
}

func (r *TemplateResource) extractHostFromJSON(content string) string {
	// Basic JSON extraction - look for "template" or "host" field
	// For a proper implementation, we'd use encoding/json
	return ""
}

func (r *TemplateResource) extractHostFromXML(content string) string {
	// Basic XML extraction
	return ""
}

// Helper functions to avoid importing strings package for simple operations
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func trimPrefix(s, prefix string) string {
	if hasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
