# Paralisus — Progress Tracker
> Last updated: 2026-05-18

## Project Summary
Class registration system built on an async **Event Bus + Worker Pool** architecture.  
Stack: Go · Redis · MongoDB · Docker Compose · Lua scripting

---

## Architecture Overview

```
HTTP Handler → IngressRouter → [CounterBus | ResponseBus | LoggingBus]
                                      ↓               ↓            ↓
                               (not impl)     RegisterWorker  LoggingWorker
                                              (Redis/Lua)     (MongoDB)
```

---

## Progress

### Infrastructure
| Item | Status |
|------|--------|
| Redis | ✅ Deployed |
| MongoDB | ✅ Deployed |
| App server | ❌ Not deployed |
| CI/CD (`.github/ci.yml`) | ⚠️ Empty — not configured |

---

### Core Modules

#### `init/` — Bootstrap
| File | Status | Notes |
|------|--------|-------|
| `init.go` — Redis + MongoDB init | ✅ Done | Pool size hardcoded to 10 |
| `settings.go` — Config loader | ⚠️ Partial | Reads `config.yml` but `LoadConfig()` hardcodes path to `./config/local.yml`; commented-out Grafana/Logger/Kafka settings |
| `run.go` — `Run()` entry | ✅ Done | |

#### `internal/event/`
| File | Status | Notes |
|------|--------|-------|
| `event.go` — Event + Request types | ✅ Done | |

#### `internal/manager/`
| File | Status | Notes |
|------|--------|-------|
| `manager.go` — IngressRouter + 3 buses | ✅ Done | CounterBus exists but has no worker |

#### `internal/worker/`
| File | Status | Notes |
|------|--------|-------|
| `worker.go` — Base Worker struct | ✅ Done | |
| `RegisterWorker/register.worker.go` | ✅ Done | Lua script runs slot decrement; drain-on-cancel implemented |
| `LoggingWorker/logging.worker.go` | ⚠️ Bug | `spawnWorker` drain loop has no `default` exit — goroutine will block forever on ctx cancel |

#### `internal/service/`
| File | Status | Notes |
|------|--------|-------|
| `warm-up.go` — WarmUpClient | ✅ Done | Both `BulkWarmup` (pipeline) and `warmup` (worker pool) implemented |
| `trackworker.go` | ❌ Empty | Placeholder only — all code commented out |

#### `internal/lua-scripting/`
| File | Status | Notes |
|------|--------|-------|
| `load-script.go` | ✅ Done | |
| `scripts/script.register.lua` | ✅ Done | |
| `scripts/script.get-class.lua` | ✅ Done | |

#### `cmd/warmup/`
| File | Status | Notes |
|------|--------|-------|
| `main.go` — Warmup CLI | ✅ Done | Supports `--csv` and `--unload` flags |

#### `main.go` — HTTP Server
| Status | Notes |
|--------|-------|
| ⚠️ Skeleton | Router created but HTTP handler is empty; buses not wired to workers; Redis/Mongo clients unused |

#### `Makefile`
| Target | Status |
|--------|--------|
| `infra-up` / `infra-down` | ✅ Done |
| `warm-up` / `unload` | ✅ Done |
| `run` | ✅ Done |
| `build` | ⚠️ Bug — `--path` is not a valid `go build` flag |

---

## To-Do List

### 🔴 Critical / Blockers

- [ ] **Wire buses to workers in `main.go`** — `IngressRouter` is created but `RegisterBus` and `LoggingBus` are never started; the HTTP handler is empty
- [ ] **Implement HTTP handler** (`/class/id=?`) — parse request, build `Event`, publish to the correct bus
- [ ] **Fix `LoggingWorker` drain loop** — the `ctx.Done()` branch has no `default` case, causing the goroutine to block indefinitely instead of draining and exiting
- [ ] **Fix `Makefile` build target** — `go build --path=` is invalid; should be `go build -o <binary> ./...` or similar

### 🟡 Incomplete / Partial

- [ ] **Implement `CounterBus` worker** — bus channel exists in `IngressRouter` but there is no worker consuming it
- [ ] **Implement `trackworker.go`** — file is a commented-out placeholder; intended to track worker state (load monitoring for scale up/down)
- [ ] **Fix `settings.go` config loading** — `viper.SetConfigFile("yml")` sets the *type*, not the extension; should use `viper.SetConfigType("yaml")` and `viper.SetConfigName("config")` to match the actual `config.yml` at root
- [ ] **`StudentID` should be an array** — `RegisterPayload.StudentID` is `string` with a comment noting it needs to be `[]string`
- [ ] **Pass Redis/Mongo clients into buses** — `main.go` discards both clients with `_, _ = rd, db`

### 🟢 Nice-to-Have / Future

- [ ] **CI/CD pipeline** — `.github/ci.yml` is empty; add build + test workflow
- [ ] **Graceful shutdown** — no `context.WithCancel` or signal handling in `main.go`; buses won't drain on SIGINT
- [ ] **Worker scale up/down** — architecture doc mentions dynamic scaling based on bus load; not yet implemented
- [ ] **Observability** — Grafana/Logger settings are commented out in config; no metrics or structured logging
- [ ] **Error handling** — several places use `panic(err)` instead of returning errors gracefully
- [ ] **Tests** — `context.md` notes testing is not currently required, but unit tests for Lua script logic and worker drain behavior would be valuable
