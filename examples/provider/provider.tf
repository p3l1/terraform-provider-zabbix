terraform {
  required_providers {
    zabbix = {
      source = "p3l1/zabbix"
    }
  }
}

provider "zabbix" {
  url       = "https://zabbix.example.com/api_jsonrpc.php"
  api_token = "your-api-token"
}