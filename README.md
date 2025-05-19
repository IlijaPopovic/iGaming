# Tournament Service

## Quick Start

Get the project up and running with just one command:

```bash
make start-project
```

---

## Project Structure Overview

```
.
├── cmd/
├── internal/
├── helperss/
├── scripts/
├── docs/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod / go.sum
```

---

## Main Components

### `cmd/main.go`

- The entry point of the application. Initializes configuration, server, and routes.

---

## internal/

### `config/`

- `config.go`: Loads environment variables and application settings.
- `db.go`: Connects to the MySQL database.

### `server/`

- `router.go`: Initializes and registers all API routes.

### `handlers/`

Responsible for HTTP request handling.

- `player_handler.go`: Player creation and retrieval logic.
- `tournament_handler.go`: Tournament creation and listing.
- `tournament_bet_handler.go`: Handles placing bets on tournaments.
- `ranking_handler.go`: Returns rankings based on player balance.
- `errors.go`: Standardized error response formatting.

#### `handlers/dtos/`

- Data Transfer Objects used for requests and responses:

  - `player.go`
  - `tournament.go`
  - `tournament_bet.go`

### `models/`

- Database models for:

  - `player.go`
  - `tournament.go`
  - `tournament_bet.go`
  - `player_rankings.go`
  - `tournament_result.go`

### `repository/`

- Responsible for DB interactions:

  - `player_repository.go`
  - `tournament_repository.go`
  - `tournament_bet_repository.go`
  - `tournament_result_repository.go`

### `migrations/`

- `001_init_schema.up.sql`: Initial SQL schema for database setup.

---

## helpers/

- `http_helpers.go`: Utility functions for handling JSON responses and errors.

---

## docs/

- Swagger auto-generated files:

  - `docs.go`: Initialization for Swagger UI.
  - `swagger.yaml` / `swagger.json`: API schema.

- Contains endpoints for:

  - Players: create, list
  - Tournaments: create, list, distribute prizes
  - Bets: place, list
  - Rankings: list

---

## scripts/

- `start.sh`: Starts the service.
- `wait_for_db.sh`: Waits for DB to become available before starting the service.

---

## Docker

- `Dockerfile`: Builds the Go application.
- `docker-compose.yml`: Spins up the app with a PostgreSQL database.

---

## API Features (Swagger)

- `GET /players` – List all players
- `POST /players` – Register a new player
- `GET /tournaments` – List all tournaments
- `POST /tournaments` – Create a new tournament
- `POST /tournaments/prizes/{id}` – Distribute prizes for a tournament
- `GET /bets` – List all bets
- `POST /bets` – Place a bet
- `GET /rankings` – Get player rankings

## Lessons Learned and Challenges

- **Swagger in Go vs C#**
  I usualy used Postman collections or the Microsoft‑style Swagger (Swashbuckle) in C#, so switching to swaggo/swag in Go was eye‑opening. The core concepts (comments → JSON schema → UI) are the same, but the annotation syntax and comment placement differ. I overcame this by reading the GoDoc examples, experimenting with a small endpoint, and iterating until the generated `swagger.yaml` matched what I expected.
- **SQL Migrations with Goose**
  Here I chose goose. Writing idempotent, reversible migration took some trial and error.
- **Waiting for Dependencies**
  Docker would often spin up the API before MySQL was ready, causing connection errors. I wrote a simple `wait_for_db.sh`, which made the start‑up process reliable in CI/CD and local Docker Compose.
- **Consistency in Error Handling**
  Ensuring every handler returned JSON in the same shape took some refactoring. I created `http_helpers.go` with helper functions like `WriteJSON(w, code, payload)` and a central `ErrorResponse` struct, so adding a new error case never meant copy‑pasting boilerplate.

---

## Efficiency & Edge Cases

Indexed & Constrained: I added indexes on player_id, tournament_id, balances, and dates. I also enforced data rules (email format, valid dates) at the schema level.

Atomic Distribution: The procedure runs inside a transaction with SELECT … FOR UPDATE on the tournament row to prevent race conditions.

Fast-Fail Guards: We immediately raise errors if there are no bets or prizes already distributed, skipping temp tables.

Set-Based Logic: Aggregations and rankings happen with temporary tables and window functions—no looping over rows.

Cleanup Safety: All temp tables are dropped at the end, and InnoDB ensures a full rollback on errors.

## TODO

Replace temporary tables in DistributePrizes with CTEs for clearer, inline data staging.

Experiment with chaining CTEs for complex prize calculations and player analytics.

Benchmark using CTEs vs temp tables to measure performance gains.

Refine the prize split logic to handle ties and variable tier distributions accurately.
