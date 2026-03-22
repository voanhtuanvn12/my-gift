# my-gift

REST + gRPC API service built with **Go**, **Iris**, **GORM**, and **PostgreSQL**.

## Tech Stack

| Layer      | Library                                                                             |
|------------|-------------------------------------------------------------------------------------|
| HTTP       | [Iris v12](https://iris-go.com)                                                     |
| gRPC       | [google.golang.org/grpc](https://grpc.io) + [grpchan](https://github.com/fullstorydev/grpchan) |
| ORM        | [GORM](https://gorm.io)                                                             |
| Database   | PostgreSQL                                                                          |
| Config     | [Viper](https://github.com/spf13/viper)                                             |
| Logger     | [Zap](https://github.com/uber-go/zap)                                               |
| API Docs   | [swaggo/swag v2](https://github.com/swaggo/swag) + [Scalar UI](https://scalar.com) |
| DI         | [Wire](https://github.com/google/wire)                                              |

## Project Structure

```
my-gift/
├── cmd/server/
│   ├── main.go          # Entry point
│   ├── app.go           # Iris + gRPC server assembly
│   ├── wire.go          # Wire injectors (build tag: wireinject)
│   └── wire_gen.go      # Wire generated (do not edit)
├── internal/
│   ├── sample/
│   │   ├── domain.go        # Interfaces + DTOs
│   │   ├── model.go         # GORM model
│   │   ├── repo.go          # PostgreSQL repository
│   │   ├── repo_dummy.go    # In-memory repository (no DB needed)
│   │   ├── service.go       # Business logic
│   │   ├── handler_http.go  # Iris MVC controller (REST)
│   │   ├── handler_grpc.go  # gRPC handler
│   │   └── provider.go      # Wire providers
│   └── infra/
│       ├── database.go  # GORM + PostgreSQL setup
│       └── logger.go    # Zap logger setup
├── pkg/
│   ├── errors/          # App error types
│   └── validator/       # Request validation
├── proto/sample/v1/     # Protobuf definitions
├── gen/proto/sample/v1/ # Generated Go code (do not edit)
├── configs/config.go    # Viper config loader
├── docs/                # Generated OAS 3.1 docs (do not edit)
├── migrations/          # SQL migrations
├── buf.yaml             # Buf config
├── buf.gen.yaml         # Buf codegen config
├── .env                 # Local env vars (git-ignored)
├── .env.example         # Env template
├── Makefile
└── Dockerfile
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL (skip if using dummy mode)
- Tool installation (first time only):

```bash
make swagger-install   # swag v2 CLI
make proto-install     # buf CLI
```

### 1. Clone & configure

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 2. Run without a database (dummy mode)

In-memory repository — no PostgreSQL required. Data resets on restart.

```bash
make run-dummy
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

make run
```

### 4. Build binary

```bash
make build
./bin/server
```

## Environment Variables

| Variable         | Default            | Description                                   |
|------------------|--------------------|-----------------------------------------------|
| `APP_NAME`       | `my-gift`          | Application name                              |
| `APP_HOST`       | `0.0.0.0`          | Bind host                                     |
| `APP_PORT`       | `8080`             | HTTP server port                              |
| `APP_GRPC_PORT`  | `50051`            | Native gRPC server port                       |
| `APP_ENV`        | `development`      | `development` / `production` / `dummy`        |
| `DB_HOST`        | `localhost`        | PostgreSQL host                               |
| `DB_PORT`        | `5432`             | PostgreSQL port                               |
| `DB_USER`        | `postgres`         | PostgreSQL user                               |
| `DB_PASSWORD`    | —                  | PostgreSQL password                           |
| `DB_NAME`        | `my_gift`          | Database name                                 |
| `DB_SSLMODE`     | `disable`          | SSL mode                                      |
| `DB_TIMEZONE`    | `Asia/Ho_Chi_Minh` | Timezone                                      |
| `LOG_LEVEL`      | `info`             | `debug` / `info` / `warn` / `error`           |

## Servers

| Server              | Port    | Protocol        | Usage                                  |
|---------------------|---------|-----------------|----------------------------------------|
| Iris HTTP           | `8080`  | HTTP/1.1+2      | REST API, Scalar docs, grpchan HTTP/1.1 |
| Native gRPC         | `50051` | HTTP/2          | grpcurl, gRPC SDK clients              |

## REST API Endpoints

| Method   | Path                    | Description         |
|----------|-------------------------|---------------------|
| `GET`    | `/health`               | Health check        |
| `GET`    | `/docs`                 | Scalar UI (OAS 3.1) |
| `GET`    | `/openapi.json`         | Raw OAS 3.1 spec    |
| `GET`    | `/api/v1/samples`       | List samples        |
| `POST`   | `/api/v1/samples`       | Create sample       |
| `GET`    | `/api/v1/samples/:id`   | Get sample by ID    |
| `PUT`    | `/api/v1/samples/:id`   | Update sample       |
| `PATCH`  | `/api/v1/samples/:id`   | Partial update      |
| `DELETE` | `/api/v1/samples/:id`   | Delete sample       |

## gRPC

### Native gRPC (port 50051)

```bash
# List all services
grpcurl -plaintext localhost:50051 list

# List methods
grpcurl -plaintext localhost:50051 list sample.v1.SampleService

# Call ListSamples
grpcurl -plaintext -d '{"page":1,"limit":10}' \
  localhost:50051 sample.v1.SampleService/ListSamples

# Call CreateSample
grpcurl -plaintext -d '{"name":"hello","description":"world"}' \
  localhost:50051 sample.v1.SampleService/CreateSample

# Call GetSample
grpcurl -plaintext -d '{"id":1}' \
  localhost:50051 sample.v1.SampleService/GetSample
```

### gRPC-over-HTTP/1.1 via grpchan (port 8080)

For clients or proxies that don't support HTTP/2.
Route format: `POST /grpc/<package>.<Service>/<Method>`

```bash
curl -X POST http://localhost:8080/grpc/sample.v1.SampleService/ListSamples \
  -H "Content-Type: application/grpc+proto"
```

## API Docs (OAS 3.1 + Scalar)

```bash
# Regenerate docs after changing annotations
make swagger

# Open Scalar UI
open http://localhost:8080/docs

# Raw OAS 3.1 spec
open http://localhost:8080/openapi.json
```

## Makefile Commands

| Command               | Description                                    |
|-----------------------|------------------------------------------------|
| `make run`            | swagger + wire + run server (with DB)          |
| `make run-dummy`      | swagger + wire + run server (in-memory, no DB) |
| `make build`          | swagger + build binary                         |
| `make swagger`        | Regenerate OAS 3.1 docs                        |
| `make swagger-install`| Install swag v2 CLI                            |
| `make wire`           | Regenerate Wire DI code                        |
| `make proto`          | Regenerate Go code from .proto files           |
| `make proto-install`  | Install buf CLI                                |
| `make tidy`           | Run `go mod tidy`                              |
| `make lint`           | Run golangci-lint                              |

## Docker

```bash
docker build -t my-gift .

# With PostgreSQL
docker run -p 8080:8080 -p 50051:50051 --env-file .env my-gift

# Dummy mode (no DB)
docker run -p 8080:8080 -p 50051:50051 -e APP_ENV=dummy my-gift
```
