## GO-DOT-SERVER

Simple HTTPS server for hosting and testing Godot Web exports locally.

It serves files from `./public`, adds the headers needed for cross-origin isolated Godot builds, and includes a basic `/healthz` endpoint for deploy checks.

## Versions

- Godot: 4.2.1
- Go: 1.23.3

## Features

- HTTPS static hosting for Godot web exports
- Security headers:
  - `Cross-Origin-Opener-Policy: same-origin`
  - `Cross-Origin-Embedder-Policy: require-corp`
- Asset caching rules for common game/static file extensions
- Gzip compression for text assets (`.html`, `.js`, `.css`, `.json`, `.txt`, `.svg`)
- Graceful shutdown on `SIGINT`/`SIGTERM`
- Health endpoint at `/healthz`
- Runtime config via environment variables

## Configuration

Defaults:

- `ADDR=0.0.0.0:8080`
- `CERT_FILE=./certs/srvr.crt`
- `KEY_FILE=./certs/priv.key`
- `PUBLIC_DIR=./public`

## Setup

1. Put your Godot web export files in `./public`.
2. Generate local certs:

```bash
mkdir -p certs
openssl req -x509 -newkey rsa:2048 -sha256 -nodes \
  -keyout certs/priv.key \
  -out certs/srvr.crt \
  -days 365 \
  -subj "/CN=localhost"
```

3. Start the server:

```bash
go run ./cmd/server
```

4. Open:

`https://127.0.0.1:8080`

## Development

Run tests:

```bash
go test ./...
```

Build binary:

```bash
go build -o build/go-dot-server ./cmd/server
```
