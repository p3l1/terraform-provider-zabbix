# Look up a host by technical name
data "zabbix_host" "server01" {
  host = "server01"
}

# Use the data source to get host information
output "server01_id" {
  value = data.zabbix_host.server01.id
}

output "server01_name" {
  value = data.zabbix_host.server01.name
}

output "server01_status" {
  value = data.zabbix_host.server01.status
}

output "server01_groups" {
  value = data.zabbix_host.server01.groups
}

output "server01_interfaces" {
  value = data.zabbix_host.server01.interfaces
}
