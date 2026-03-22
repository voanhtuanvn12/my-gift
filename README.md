# my-gift

RESTful API service built with **Go**, **Iris**, **GORM**, and **PostgreSQL**.

## Tech Stack

| Layer      | Library                          |
|------------|----------------------------------|
| HTTP       | [Iris v12](https://iris-go.com)  |
| ORM        | [GORM](https://gorm.io)          |
| Database   | PostgreSQL                       |
| Config     | [Viper](https://github.com/spf13/viper) |
| Logger     | [Zap](https://github.com/uber-go/zap) |
| API Docs   | [swaggo/swag v2](https://github.com/swaggo/swag) + [Scalar UI](https://scalar.com) |
| DI         | [Wire](https://github.com/google/wire) |

## Project Structure

```
my-gift/
├── cmd/server/
│   ├── main.go          # Entry point
│   └── wire.go          # Wire DI providers (codegen)
├── internal/
│   ├── sample/
│   │   ├── domain.go        # Interfaces + DTOs
│   │   ├── model.go         # GORM model
│   │   ├── repo.go          # PostgreSQL repository
│   │   ├── repo_dummy.go    # In-memory repository (testing)
│   │   ├── service.go       # Business logic
│   │   ├── handler_http.go  # Iris MVC controller
│   │   └── provider.go      # Wire providers
│   └── infra/
│       ├── database.go  # GORM + PostgreSQL setup
│       └── logger.go    # Zap logger setup
├── pkg/
│   ├── errors/          # App error types
│   └── validator/       # Request validation
├── configs/config.go    # Viper config loader
├── docs/                # Generated Swagger docs (do not edit)
├── migrations/          # SQL migrations
├── .env                 # Local env vars (git-ignored)
├── .env.example         # Env template
├── Makefile
└── Dockerfile
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL (skip if using dummy mode)
- `swag` CLI — install once: `make swagger-install`

### 1. Clone & configure

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 2. Run without a database (dummy mode)

Dùng in-memory repository, không cần PostgreSQL. Data mất khi restart.

```bash
APP_ENV=dummy go run ./cmd/server/...
```

### 3. Run with PostgreSQL

```bash
# Start PostgreSQL (Docker)
docker run -d \
  --name my-gift-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=my_gift \
  -p 5432:5432 \
  postgres:16-alpine

# Run server
make run
```

### 4. Build binary

```bash
make build
./bin/server
```

## Environment Variables

| Variable       | Default         | Description              |
|----------------|-----------------|--------------------------|
| `APP_NAME`     | `my-gift`       | Application name         |
| `APP_HOST`     | `0.0.0.0`       | Bind host                |
| `APP_PORT`     | `8080`          | Bind port                |
| `APP_ENV`      | `development`   | `development` / `production` / `dummy` |
| `DB_HOST`      | `localhost`     | PostgreSQL host          |
| `DB_PORT`      | `5432`          | PostgreSQL port          |
| `DB_USER`      | `postgres`      | PostgreSQL user          |
| `DB_PASSWORD`  | —               | PostgreSQL password      |
| `DB_NAME`      | `my_gift`       | Database name            |
| `DB_SSLMODE`   | `disable`       | SSL mode                 |
| `DB_TIMEZONE`  | `Asia/Ho_Chi_Minh` | Timezone              |
| `LOG_LEVEL`    | `info`          | `debug` / `info` / `warn` / `error` |

## API Endpoints

| Method   | Path                   | Description         |
|----------|------------------------|---------------------|
| `GET`    | `/health`              | Health check        |
| `GET`    | `/docs`                | Scalar UI (OAS 3.1) |
| `GET`    | `/openapi.json`        | Raw OAS 3.1 spec    |
| `GET`    | `/api/v1/samples`      | List samples        |
| `POST`   | `/api/v1/samples`      | Create sample       |
| `GET`    | `/api/v1/samples/:id`  | Get sample by ID    |
| `PUT`    | `/api/v1/samples/:id`  | Update sample       |
| `PATCH`  | `/api/v1/samples/:id`  | Partial update      |
| `DELETE` | `/api/v1/samples/:id`  | Delete sample       |

## API Docs (OAS 3.1 + Scalar)

```bash
# Install swag v2 CLI (first time only)
make swagger-install

# Regenerate docs after changing annotations
make swagger

# Open Scalar UI in browser
open http://localhost:8080/docs

# Raw OAS 3.1 spec
open http://localhost:8080/openapi.json
```

## Makefile Commands

| Command               | Description                          |
|-----------------------|--------------------------------------|
| `make run`            | Generate swagger docs + run server   |
| `make build`          | Generate swagger docs + build binary |
| `make swagger`        | Regenerate Swagger docs              |
| `make swagger-install`| Install swag CLI                     |
| `make wire`           | Regenerate Wire DI code              |
| `make tidy`           | Run `go mod tidy`                    |
| `make lint`           | Run golangci-lint                    |

## Docker

```bash
docker build -t my-gift .
docker run -p 8080:8080 --env-file .env my-gift
```
