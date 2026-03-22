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
│   ├── middleware/
│   │   ├── middleware.go    # WrapRouter / UseRouter / UseGlobal / Use / UseError / Done / DoneGlobal
│   │   └── jwt.go           # JWTVerify / GetClaims / GenerateToken
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
| `JWT_SECRET`   | `change-me-in-production` | HMAC-SHA256 signing key  |
| `JWT_EXPIRY`   | `24h`           | Token expiry (Go duration string)    |

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

## Implementation Steps

Thứ tự xây dựng project từ đầu (để tham khảo khi tạo domain mới):

### Bước 1 — Cấu trúc nền (`configs`, `infra`)
1. `configs/config.go` — load env vars với Viper, định nghĩa `Config` struct
2. `internal/infra/logger.go` — khởi tạo Zap logger theo env
3. `internal/infra/database.go` — kết nối PostgreSQL qua GORM

### Bước 2 — Domain layer (`internal/<domain>/`)
4. `domain.go` — định nghĩa entity, DTOs (Request/Response), interface `Service` và `Repository`
5. `model.go` — GORM model + mapper `ToDomain()` / `fromDomain()`
6. `repo.go` — implement `Repository` với GORM (PostgreSQL)
7. `repo_dummy.go` — implement `Repository` in-memory (dùng khi `APP_ENV=dummy`)
8. `service.go` — implement `Service`, chứa business logic
9. `handler_http.go` — Iris MVC controller, map HTTP ↔ Service
10. `provider.go` — Wire provider functions (`ProvideRepository`, `ProvideService`, `ProvideController`)

### Bước 3 — Middleware (`internal/middleware/`)
11. `middleware.go` — 7 lớp middleware theo thứ tự Iris pipeline:
    - `WrapRouter` → low-level nhất, CORS, rate limit
    - `UseRouter` → request ID, access log
    - `UseGlobal` → chạy cho mọi route kể cả error pages
    - `Use` → chạy cho route thường (auth, business logic)
    - `UseError` → chỉ chạy cho error handler
    - `Done` → cleanup sau route thường
    - `DoneGlobal` → cleanup cho mọi route
12. `jwt.go` — `JWTVerify`, `GetClaims`, `GenerateToken`

### Bước 4 — Wiring (`cmd/server/`)
13. `wire.go` — khai báo Wire provider sets và injector functions
14. `wire_gen.go` — **auto-generated** bởi `make wire`, không sửa tay
15. `app.go` — khởi tạo Iris app, đăng ký middleware, routes, MVC error handler
16. `main.go` — entry point, chọn `InitializeApp` hay `InitializeAppDummy` theo env

### Bước 5 — API Docs
17. Thêm Swaggo annotations vào `main.go` (global) và `handler_http.go` (per-endpoint)
18. `make swagger` → sinh `docs/`

### Quy tắc khi thêm domain mới

```
internal/<domain>/
├── domain.go       # entity + DTOs + interface
├── model.go        # GORM model
├── repo.go         # PostgreSQL impl
├── repo_dummy.go   # in-memory impl
├── service.go      # business logic
├── handler_http.go # HTTP controller
└── provider.go     # Wire providers
```

Sau đó:
- Thêm provider set vào `wire.go`
- Đăng ký route trong `app.go`
- Chạy `make wire` để tái sinh `wire_gen.go`
- Chạy `make swagger` để cập nhật docs

## Docker

```bash
docker build -t my-gift .
docker run -p 8080:8080 --env-file .env my-gift
```
