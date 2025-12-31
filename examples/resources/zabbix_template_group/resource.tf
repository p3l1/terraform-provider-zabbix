# Create a template group for application templates
resource "zabbix_template_group" "applications" {
  name = "Templates/Applications"
}

# Create a template group for OS templates
resource "zabbix_template_group" "operating_systems" {
  name = "Templates/Operating systems"
}
