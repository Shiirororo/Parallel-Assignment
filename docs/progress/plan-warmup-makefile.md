# Plan: Deploy Warm-up Stage via Makefile

## Goal
Expose the warm-up stage as a standalone Makefile target so an admin can pre-load `class.csv` into Redis without starting the full HTTP server.

---

## What the warm-up does
`internal/service/warm-up.go` reads a CSV (`class.csv`) and writes each row as a Redis key (`class_id → slot`) using either:
- `BulkWarmup` — single pipeline, best for large batches
- `warmup` — worker-pool approach, configurable concurrency

The entry point needs to be wired into a runnable `cmd/` binary.

---

## Steps

### 1. Create `cmd/warmup/main.go`
A minimal CLI entry point that:
- Calls `bootstrap.Run()` to get the Redis client (reads `config.yml`)
- Instantiates `WarmUpClient`
- Calls `BulkWarmup(ctx, "class.csv")`
- Exits cleanly

```
cmd/
└── warmup/
    └── main.go
```

### 2. Export `WarmUpClient` constructor and `BulkWarmup`
`newWarmUpClient` is currently unexported. Export it as `NewWarmUpClient` so `cmd/warmup` can use it.

### 3. Add Makefile targets

```makefile
WARMUP_CSV ?= class.csv

warmup:
	go run ./cmd/warmup --csv=$(WARMUP_CSV)

unload:
	go run ./cmd/warmup --unload

infra-up:
	docker compose up -d redis mongo

infra-down:
	docker compose down
```

Typical admin workflow:
```
make infra-up      # start Redis + MongoDB
make warmup        # load class.csv into Redis
make run           # start the HTTP server
```

### 4. Wire CLI flags in `cmd/warmup/main.go`
Accept two flags:
- `--csv=<path>` (default: `class.csv`) — path to the CSV file
- `--unload` — flush the Redis DB instead of loading

---

## Acceptance Criteria
- `make warmup` loads all rows from `class.csv` into Redis and prints `Bulk import complete`
- `make warmup WARMUP_CSV=other.csv` works with a custom file
- `make unload` flushes Redis and prints `Data flushed`
- Neither target starts the HTTP server
