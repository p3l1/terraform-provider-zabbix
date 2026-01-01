# Create a simple template with metadata
resource "zabbix_template_group" "custom" {
  name = "Custom Templates"
}

resource "zabbix_template" "example" {
  host        = "my_custom_template"
  name        = "My Custom Template"
  description = "A custom monitoring template"
  groups      = [zabbix_template_group.custom.id]

  tags = [
    {
      tag   = "environment"
      value = "production"
    },
    {
      tag   = "team"
      value = "platform"
    }
  ]
}

# Import an official Zabbix template from YAML content
# The template will be created with its embedded metadata including groups
resource "zabbix_template" "apache" {
  source_format  = "yaml"
  source_content = file("apache_template.yaml")
}
