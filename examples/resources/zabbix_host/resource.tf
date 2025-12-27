# Create a basic host with an agent interface
resource "zabbix_host" "server01" {
  host   = "server01"
  name   = "Production Server 01"
  groups = [zabbix_host_group.linux.id]
  status = 0

  interfaces = [{
    type   = "agent"
    ip     = "192.168.1.100"
    dns    = ""
    port   = "10050"
    main   = true
    use_ip = true
  }]
}

# Create a host with multiple interfaces and tags
resource "zabbix_host" "server02" {
  host   = "server02"
  name   = "Production Server 02"
  groups = [zabbix_host_group.linux.id, zabbix_host_group.web.id]
  status = 0

  interfaces = [
    {
      type   = "agent"
      ip     = "192.168.1.101"
      dns    = "server02.example.com"
      port   = "10050"
      main   = true
      use_ip = true
    },
    {
      type   = "snmp"
      ip     = "192.168.1.101"
      dns    = ""
      port   = "161"
      main   = true
      use_ip = true
    }
  ]

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

# Create a host with linked templates
resource "zabbix_host" "webserver" {
  host      = "webserver01"
  name      = "Web Server 01"
  groups    = [zabbix_host_group.web.id]
  templates = ["10001"] # Template ID for Linux by Zabbix agent

  interfaces = [{
    type   = "agent"
    ip     = "192.168.1.200"
    dns    = ""
    port   = "10050"
    main   = true
    use_ip = true
  }]
}

# Create a disabled host
resource "zabbix_host" "maintenance" {
  host   = "maintenance-server"
  name   = "Server Under Maintenance"
  groups = [zabbix_host_group.linux.id]
  status = 1 # disabled

  interfaces = [{
    type   = "agent"
    ip     = "192.168.1.150"
    dns    = ""
    port   = "10050"
    main   = true
    use_ip = true
  }]
}

# Referenced host groups
resource "zabbix_host_group" "linux" {
  name = "Linux servers"
}

resource "zabbix_host_group" "web" {
  name = "Web servers"
}
