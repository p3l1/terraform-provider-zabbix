# Look up an existing host group by name
data "zabbix_host_group" "linux" {
  name = "Linux servers"
}

# Use the host group ID in other resources
output "linux_group_id" {
  value = data.zabbix_host_group.linux.id
}
