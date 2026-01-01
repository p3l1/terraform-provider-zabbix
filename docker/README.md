# Zabbix Test Environment

Local Zabbix stack for running acceptance tests.

## Quick Start

```bash
# Start the stack
docker compose up -d

# Wait for services to be healthy
docker compose ps

# Stop the stack
docker compose down

# Stop and remove all data
docker compose down -v
```

## Services

| Service | Port | Description |
|---------|------|-------------|
| postgres | 5432 (internal) | PostgreSQL 16 database |
| zabbix-server | 10051 (internal) | Zabbix Server |
| zabbix-web | 8080 | Web UI and API endpoint |

## API Access

- **URL**: `http://localhost:8080/api_jsonrpc.php`
- **Default credentials**: Admin / zabbix
- **Pre-configured API token**:

```
071fb9d2e8f72cf9c40128f0f5aab3def1bab0893413314b083fdcb4551eb01a
```

The API token is automatically created by the `db-init` service when starting fresh.

## Environment Variables for Tests

```bash
export ZABBIX_URL="http://localhost:8080/api_jsonrpc.php"
export ZABBIX_API_TOKEN="071fb9d2e8f72cf9c40128f0f5aab3def1bab0893413314b083fdcb4551eb01a"
```
