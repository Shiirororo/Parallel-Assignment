# Paralisus

A high-throughput class registration service built in Go. Uses an event-driven worker bus architecture to handle concurrent registrations with Redis for slot management and MongoDB for persistence.

## Architecture

Incoming requests are dispatched to three parallel worker buses (4 workers each):

- **RegisterBus** — decrements remaining slots in Redis via Lua scripting (atomic)
- **CounterBus** — tracks active registration counts
- **LoggingBus** — persists registration records to MongoDB

A local in-memory cache layer syncs from Redis every 5 seconds to reduce read pressure.

## Prerequisites

- Go 1.25+
- Docker & Docker Compose

## Getting Started

**1. Start infrastructure**
```bash
make infra-up
```

**2. Seed class data**
```bash
make warm-up
```

**3. Run the server**
```bash
make run
```

Server listens on `:36789`.

## API

| Method | Endpoint                  | Description                        |
|--------|---------------------------|------------------------------------|
| GET    | `/api/class/getClassInfo` | Get remaining slots for class IDs  |
| POST   | `/api/class/register`     | Register a student for a class     |
| POST   | `/api/class/unregister`   | Unregister a student from a class  |

See [`docs/instruction/API.md`](docs/instruction/API.md) for full request/response details.

## Configuration

Edit `config.yml` to configure Redis, MongoDB, and server settings.

```yaml
server:
  port: 36789

redis:
  host: 127.0.0.1
  port: 6379

mongo:
  host: 127.0.0.1
  port: 27017
  database: paralisus
```

## Build

```bash
make build   # produces ./paralisus binary
```
