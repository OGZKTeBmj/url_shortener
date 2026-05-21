# URL Shortener

A high-performance URL shortening service built with Go. The project is fully containerized using Docker and managed via a streamlined Makefile workflow.

## Prerequisites

Ensure you have the following dependencies installed locally:
* Docker & Docker Compose
* GNU Make
* Go (v1.26+) — *Optional: required only for building and running outside of Docker*

---

## Quick Start

### 1. Environment Setup
Clone the configuration template and create your local environment file:
```bash
cp .env.example .env
```
Note: Open the newly created `.env` file and update the default credentials for PostgreSQL and Redis if necessary.

### 2. Spin Up the Stack
Build the application images, provision PostgreSQL and Redis instances, and automatically run database migrations:
```bash
make up
```
To stream live logs from all running containers:
```bash
make logs
```

### 3. Tear Down the Stack
To stop and remove all application containers:
```bash
make down
```

---

## Infrastructure Management (Docker)

Isolated commands to manage individual database and caching layers during development:

### PostgreSQL
* `make db-up` — Spin up the PostgreSQL container in detached mode.
* `make db-down` — Stop the PostgreSQL container without destroying data.

### Redis
* `make redis-up` — Spin up the Redis container in detached mode.
* `make redis-down` — Stop the Redis container.
* `make redis-logs` — Tail live logs from the Redis container.
* `make redis-cli` — Drop into an interactive Redis CLI session (authentication is handled automatically via `.env`).

### Deep Clean
To wipe the entire stack, drop all persistent database volumes, and rebuild from scratch:
```bash
make reset
```

---

## Database Migrations

Schema evolution is managed using `golang-migrate` executed inside a transient Docker container.

* `make migrate-up` — Apply all pending schema upgrades.
* `make migrate-down` — Roll back the last applied migration step.
* `make migrate-reset` — Forcefully drop all existing tables and schema history.
* `make migrate-version` — Output the current migration version database state.
* `make migrate-create name=migration_name` — Generate a new pair of sequential SQL migration files (`up`/`down`) inside `./migrations`.

---

## Local Development (Native Execution)

To run the Go binary directly on your host machine while utilizing Dockerized backing services:

1. Launch the storage infrastructure:
   ```bash
   make db-up redis-up
   ```
2. Build and execute the Go binary locally:
   ```bash
   make run
   ```
3. Remove compiled binaries and clean the workspace:
   ```bash
   make clean
   ```