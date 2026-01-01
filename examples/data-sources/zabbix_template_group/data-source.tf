# Look up an existing template group by name
data "zabbix_template_group" "applications" {
  name = "Templates/Applications"
}

# Use the template group ID in other resources
output "template_group_id" {
  value = data.zabbix_template_group.applications.id
}
