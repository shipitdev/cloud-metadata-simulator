# Cloud Metadata Server Simulator (IMDS)

A lightweight Go HTTP server that simulates cloud provider Instance Metadata Service (IMDSv1) and OpenStack metadata APIs. It reads node metadata and `#cloud-config` user-data from a YAML file to serve cloud instance initialization requests.

## Features
- **IMDS Endpoints:** Standard `/latest/meta-data/` routes (`instance-id`, `hostname`, `local-ipv4`, `public-ipv4`, `public-keys/`).
- **User Data Endpoint:** `/latest/user-data` serving raw `#cloud-config` content.
- **JSON Output:** `/latest/meta-data.json` returning structured metadata.
- **Configurable:** Load metadata definitions via `--config` flag.

## Usage

### Run server
```bash
go run main.go --config metadata.yaml --port 8080
```

### Test endpoints
```bash
# Get Instance ID
curl http://localhost:8080/latest/meta-data/instance-id

# Get Public SSH Key
curl http://localhost:8080/latest/meta-data/public-keys/0/openssh-key

# Get User-Data
curl http://localhost:8080/latest/user-data
```

## Running Tests
```bash
go test -v ./...
```

## License
MIT
