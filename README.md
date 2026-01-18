# balancr

Lightweight load balancer built from scratch in Go — no external libraries for networking.

## Features

- **TCP Proxying** - Forward connections at the transport layer
- **Weighted Round-Robin** - Distribute traffic based on backend weights
- **Health Checks** - Automatic detection of dead/alive backends
- **Structured Logging** - Timestamped logs with levels (INFO, WARN, ERROR)
- **Metrics Tracking** - Request counts per backend, error tracking
- **CLI Configuration** - Configurable via command-line flags

## Architecture

```
Client Request
      |
      v
+-------------+
|   balancr   |  <-- distributes load (weighted round-robin)
+-------------+
      |
   +--+--+--+
   |  |  |  |
   v  v  v  v
+--+ +--+ +--+
|S1| |S2| |S3|  <-- backend servers
+--+ +--+ +--+
```

## Usage

```bash
# Default (port 7000, config.yaml)
go run cmd/balancr/main.go

# Custom configuration
go run cmd/balancr/main.go -port 8080 -config prod.yaml -health-interval 5 -metrics-interval 30

# View all options
go run cmd/balancr/main.go -help
```

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | 7000 | Port to listen on |
| `-config` | config.yaml | Path to config file |
| `-health-interval` | 10 | Health check interval (seconds) |
| `-metrics-interval` | 10 | Metrics output interval (seconds) |

## Configuration

```yaml
# config.yaml
backends:
  - url: "localhost:9999"
    weight: 3
  - url: "localhost:9998"
    weight: 1
  - url: "localhost:9997"
    weight: 2
```

Weights determine traffic distribution. With weights 3:1:2, over 6 requests:
- `localhost:9999` receives 3 requests
- `localhost:9998` receives 1 request
- `localhost:9997` receives 2 requests

## Project Structure

```
balancr/
├── cmd/balancr/main.go       # Entry point, CLI flags
├── internal/
│   ├── config/config.go      # YAML config parsing
│   ├── pool/pool.go          # Server pool, weighted round-robin, health checks
│   ├── proxy/proxy.go        # TCP proxying
│   ├── logger/logger.go      # Structured logging
│   └── metrics/metrics.go    # Request/error tracking
└── config.yaml               # Backend configuration
```

## Sample Output

```
2026-01-18 21:29:13 [INFO] Listening to port 7000
2026-01-18 21:29:15 [INFO] Request forwarded to localhost:9999
2026-01-18 21:29:16 [INFO] Request forwarded to localhost:9998
2026-01-18 21:29:16 [INFO] Request forwarded to localhost:9997
2026-01-18 21:29:23 [INFO] Total: 6 | Errors: 0
localhost:9999: 3
localhost:9998: 1
localhost:9997: 2
```

## Why

Built to understand load balancing internals at the TCP level, not to replace production solutions like nginx or HAProxy.

## Author

Haseeb Ahmed - [GitHub](https://github.com/haseebahmed248)
