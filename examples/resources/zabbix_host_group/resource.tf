# Create a host group for Linux servers
resource "zabbix_host_group" "linux" {
  name = "Linux servers"
}

# Create a host group for web servers
resource "zabbix_host_group" "web" {
  name = "Web servers"
}
