# balancr

Lightweight load balancer built from scratch in Go — no external libraries.

## Goal

Understand load balancing at the network level by implementing:
- TCP connection proxying
- Round-robin distribution
- Health checks
- Weighted balancing
- Request logging

## Status

**Work in progress**

- [ ] Basic proxy (forward to single backend)
- [ ] Multiple backends + round-robin
- [ ] Health checks (detect dead servers)
- [ ] Config file + weighted balancing
- [ ] Logging + metrics
- [ ] CLI flags + polish

## Architecture

```
Client Request
      |
      v
+-------------+
|   balancr   |  <-- distributes load
+-------------+
      |
   +--+--+
   |     |
   v     v
+----+ +----+
| S1 | | S2 |  <-- backend servers
+----+ +----+
```

## Run

```bash
./balancr --config config.yaml --port 8080
```

## Why

Built to understand load balancing internals, not to replace production solutions like nginx or HAProxy.
