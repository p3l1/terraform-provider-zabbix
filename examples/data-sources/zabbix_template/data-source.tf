# Look up an existing template by technical name
data "zabbix_template" "linux" {
  host = "Linux by Zabbix agent"
}

output "template_id" {
  value = data.zabbix_template.linux.id
}

output "template_uuid" {
  value = data.zabbix_template.linux.uuid
}

# The exported_content can be used for backup or drift detection
output "template_content" {
  value     = data.zabbix_template.linux.exported_content
  sensitive = true
}
