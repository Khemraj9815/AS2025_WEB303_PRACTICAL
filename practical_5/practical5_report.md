# Practical 5 — Refactoring Monolith to Microservices

Date: 2025-11-09

## Summary

This practical demonstrates migrating a simple Student Cafe application from a monolithic design to a microservices architecture. The reference implementation splits the app into:

- user-service (users)
- menu-service (menu items)
- order-service (orders)
- api-gateway (single entry point)

I exercised the code, fixed build and Docker issues, and tested the API gateway container locally.

## Takeaways (in simple words)

- Microservices split responsibilities: each service owns its models and database.
- Keep module paths consistent: use full module paths for Go modules to avoid import problems.
- Docker Compose makes it easy to run many services together, but each service must have a valid Dockerfile and correct configuration.
- Small integration points (like the order-service calling user-service and menu-service) are easy to implement with HTTP but need careful error handling and timeouts in real systems.
- Auto-migration via GORM simplifies DB setup for local testing.

## Issues faced (what I ran into) — plain language

1. Invalid docker-compose YAML in `api-gateway/docker-compose.yml`:
   - The file omitted the top-level `services:` key, causing Docker Compose to reject it.
   - Fix: added proper `services:` structure and removed `depends_on` pointing to services that were not declared in that file.

2. Missing Dockerfile for `api-gateway`:
   - Running `docker compose up` failed because there was no Dockerfile in the `api-gateway` folder.
   - Fix: added a simple Dockerfile that builds the gateway binary.

3. Go module import issues across services:
   - Some services used short module names like `order-service/database` which confuse the Go toolchain outside a special GOPATH layout.
   - Fix: for `order-service` I updated `go.mod` to a fully-qualified module path and adjusted imports. (Recommend doing the same for other services for consistent local builds.)

4. Builder Go version mismatch:
   - Docker build used `golang:1.23` but the project required `go >= 1.25.0`, causing `go mod download` to fail.
   - Fix: updated the Dockerfile to use `golang:1.25-alpine`.

5. Corrupted or duplicated source content (during edits):
   - While working I found and fixed a few files that had accidental duplicated lines or stray tokens causing compile errors.

## What I changed (short list)

- Added `api-gateway/Dockerfile` and adjusted it to Go 1.25.
- Fixed `api-gateway/docker-compose.yml` structure so it validates.
- Implemented/cleaned `order-service` models, DB connect, and `main.go` and updated module path for `order-service`.

## Quick commands I used (you can copy/paste)

Validate a compose file:
```
cd practical_5/api-gateway
docker compose config
```

Build and run the gateway locally:
```
cd practical_5/api-gateway
docker compose up --build
```

Run all services (recommended): create or use a root `docker-compose.yml` at the project root that declares all services and DB containers, then:
```
cd practical_5
docker compose up --build
```

## Suggestions / next steps

- Normalize module names for all services to fully-qualified module paths (e.g. `github.com/<org>/<repo>/practicals/practical5/<service>`) and update imports. This makes `go build` and `go mod tidy` behave predictably.
- Add a root `docker-compose.yml` that defines all services and the Postgres databases so the entire system can be started with one command.
- Add simple integration tests or a health-check endpoint to each service so the API gateway can route traffic only when a service is ready.
- Add timeouts and better error handling for inter-service HTTP calls in `order-service`.

---

This report is written in simple words to capture the main lessons and the concrete issues I fixed while doing Practical 5.
