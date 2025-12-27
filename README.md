# Terraform Provider for Zabbix

[![Tests](https://github.com/p3l1/terraform-provider-zabbix/actions/workflows/test.yml/badge.svg)](https://github.com/p3l1/terraform-provider-zabbix/actions/workflows/test.yml)

> **Status:** Early development - not yet ready for production use

A Terraform/OpenTofu provider for managing [Zabbix](https://www.zabbix.com/) monitoring infrastructure as code. Define hosts, templates, triggers, and other Zabbix resources declaratively and version-control your monitoring configuration.

## Requirements

- [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.10.7
- [Go](https://golang.org/doc/install) >= 1.25 (for building)
- [Zabbix](https://www.zabbix.com/) Server >= 7.0

## Supported Resources

### Implemented

| Resource            | Data Source         | Description        |
| ------------------- | ------------------- | ------------------ |
| `zabbix_host`       | `zabbix_host`       | Manage hosts       |
| `zabbix_host_group` | `zabbix_host_group` | Manage host groups |

### Planned

- Host interfaces
- Templates and template groups
- Items, triggers, and graphs
- Discovery rules
- Actions and media types
- Users, user groups, and roles
- Proxies and macros
- Services and SLAs

## Documentation

This provider maps Terraform resources to Zabbix API objects. For details on the underlying API, see the [Zabbix API Reference](https://www.zabbix.com/documentation/7.0/en/manual/api/reference).

## License

MIT License
