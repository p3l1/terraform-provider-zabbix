# Zabbix Terraform Provider

Terraform provider for managing Zabbix monitoring infrastructure.

## Build Commands

```bash
make build      # Build the provider
make test       # Run unit tests
make testacc    # Run acceptance tests (requires TF_ACC=1)
make lint       # Run golangci-lint
make fmt        # Format code with gofmt
make generate   # Run go generate
make install    # Build and install locally
```

## Project Structure

```
.
├── main.go                      # Provider entry point
├── internal/provider/
│   ├── provider.go              # Provider implementation
│   └── provider_test.go         # Provider tests
├── tools/tools.go               # Build tool dependencies
├── Makefile                     # Build targets
├── .golangci.yml                # Linter configuration
└── .goreleaser.yml              # Release configuration
```

## Provider Configuration

The provider supports two configuration methods:

1. **HCL Configuration:**
```hcl
provider "zabbix" {
  url       = "https://zabbix.example.com/api_jsonrpc.php"
  api_token = "your-api-token"
}
```

2. **Environment Variables:**
```bash
export ZABBIX_URL="https://zabbix.example.com/api_jsonrpc.php"
export ZABBIX_API_TOKEN="your-api-token"
```

HCL configuration takes precedence over environment variables.

## Development Workflow

1. Create a feature branch from an issue: `<issue-number>-short-description`
2. Write tests first (TDD)
3. Implement the feature
4. Run `make test` and `make build`
5. Commit with descriptive messages
6. Create PR targeting `main`
7. Merge using merge, no squashing

## Testing

### Docker Test Environment

A Docker-based Zabbix test environment is available in `docker/`:

```bash
cd docker
docker compose up -d
```

This starts Zabbix server, web frontend, and PostgreSQL database for local testing with preconfigured static API token. The static token (`071fb9d2e8f72cf9c40128f0f5aab3def1bab0893413314b083fdcb4551eb01a`) is used as a default in integration tests.

## Code Conventions

- All files must start with two-line `ABOUTME:` comments explaining the file's purpose
- Use terraform-plugin-framework (not the older SDK)
- Provider attributes that can be set via environment variables should be `Optional: true`
- Sensitive attributes must have `Sensitive: true`
