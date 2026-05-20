

# Main context
This project deploy a class register system using async Event Bus + Worker pools


This project was created to testing concurent programming. We like to apply parallel computing strategies: Worker pool, Event bus for async request progress

Current Project tree:
[4.0K]  .
├── [ 11K]  class.csv
├── [4.0K]  cmd
│   ├── [4.0K]  load-test
│   │   └── [ 498]  load-test.js
│   └── [4.0K]  warmup
│       ├── [ 12M]  main
│       ├── [ 635]  main.go
│       └── [ 489]  mock-data.py
├── [ 503]  compose.yml
├── [ 175]  config.yml
├── [1.9K]  context.md
├── [4.0K]  docs
│   ├── [4.0K]  instruction
│   │   └── [1.7K]  API.md
│   └── [4.0K]  progress
├── [1.4K]  go.mod
├── [8.1K]  go.sum
├── [4.0K]  init
│   ├── [ 889]  init.go
│   ├── [ 236]  run.go
│   └── [1.4K]  settings.go
├── [4.0K]  internal
│   ├── [4.0K]  event
│   │   └── [ 179]  event.go
│   ├── [4.0K]  lua-scripting
│   │   ├── [ 534]  load-script.go
│   │   └── [4.0K]  scripts
│   │       ├── [ 205]  script.get-class.lua
│   │       ├── [ 126]  script.register.lua
│   │       └── [ 143]  script.unregister.lua
│   ├── [4.0K]  manager
│   │   └── [ 999]  manager.go
│   ├── [4.0K]  service
│   │   ├── [ 465]  trackworker.go
│   │   └── [1.6K]  warm-up.go
│   └── [4.0K]  worker
│       ├── [4.0K]  CounterWorker
│       │   └── [ 133]  counter.worker.go
│       ├── [4.0K]  LoggingWorker
│       │   └── [1.9K]  logging.worker.go
│       ├── [4.0K]  RegisterWorker
│       │   └── [3.1K]  register.worker.go
│       └── [ 107]  worker.go
├── [1.8K]  main.go
├── [ 360]  Makefile
├── [1.8K]  plan-warmup-makefile.md
└── [4.9K]  PROGRESS.md

18 directories, 30 files

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