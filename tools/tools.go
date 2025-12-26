// ABOUTME: Build tool dependencies for the Zabbix Terraform provider.
// ABOUTME: Declares tools needed for development (linting, docs generation) as Go dependencies.

//go:build tools

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
