

# Main context
This project deploy a class register system using async Event Bus + Worker pools


This project was created to testing concurent programming. We like to apply parallel computing strategies: Worker pool, Event bus for async request progress

Current Project tree:
.
├── class.csv
├── cmd
├── compose.yml
├── config.yml
├── context.md
├── go.mod
├── go.sum
├── init
│   ├── init.go
│   ├── run.go
│   └── settings.go
├── internal
│   ├── event
│   │   └── event.go
│   ├── lua-scripting
│   │   └── script.register.lua
│   ├── manager
│   │   ├── event.go
│   │   └── manager.go
│   ├── service
│   │   ├── trackworker.go
│   │   └── warm-up.go
│   └── worker
│       ├── LoggingWorker
│       │   └── updateDB.worker.go
│       ├── ResponseWorker
│       │   └── response.worker.go
│       └── worker.go
├── main.go
└── Makefile

# Event Bus + Worker pools architect


This architect include:
    - A manager IngressRouter manage 3 Bus: ResponseBus, LoggingBus, CounterBus
    - Each bus manage its own worker and their jobs, bus allow scale up/down worker base on their load
    - The worker must finish all their job before die
    - Worker and Bus communicate through channel

There were 2 main high-concurrent jobs including:
    - Response Worker: Update remaining slot after class slot confirmation. Require high consistency
    - Logging Worker: save register status into MongoDB for later query job

# Function

## Warm-up
This allow admin pre-load class information and slot.
From a CSV file contain information and load all into Redis Cache
## Bus

## Worker
# Current State

- Infrastructure:
 - Redis: Deploy
 - MongoDB: Deploy
 - App: Not deploy


Current state not require testing