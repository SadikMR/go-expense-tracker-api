# Go Expense Tracker API

A small REST API for tracking expenses with Go, Beego, and CSV-backed persistence.

## Project overview

This service supports:
- user registration and login
- expense CRUD operations
- expense summary reporting
- Swagger API documentation
- CSV storage for users and expenses

It is built with Go modules and Beego, and it uses configuration values from `conf/app.conf`.

## Features

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/health`
- `GET /api/v1/expenses`
- `POST /api/v1/expenses`
- `GET /api/v1/expenses/:id`
- `PUT /api/v1/expenses/:id`
- `DELETE /api/v1/expenses/:id`
- `GET /api/v1/expenses/summary`

## Requirements

- Go 1.26 or newer
- Git

Optional:
- `bee` CLI for development (`go install github.com/beego/bee/v2@latest`)
- `swag` CLI for regenerating Swagger docs (`go install github.com/swaggo/swag/cmd/swag@latest`)

## Install

```bash
git clone https://github.com/SadikMR/go-expense-tracker-api.git
cd go-expense-tracker-api
go mod download
```

## Configuration

The application loads configuration from `conf/app.conf` and supports environment overrides.

Important settings:

- `httpport` — API port
- `runmode` — `dev` or `prod`
- `data_dir` — fallback directory for CSV files
- `users_csv_path` — optional explicit users CSV path
- `expenses_csv_path` — optional explicit expenses CSV path
- `log_dir` — optional log directory for production log files

Copy `conf/app.conf.example` to `conf/app.conf` and provide the proper values before running the app:

```bash
cp conf/app.conf.example conf/app.conf
```

Environment variables that can override config values:

- `PORT`
- `RUN_MODE`
- `DATA_DIR`
- `USERS_CSV_PATH`
- `EXPENSES_CSV_PATH`

## Run the application

Start the server with Go:

```bash
go run main.go
```

Or use Bee CLI if installed:

```bash
bee run
```

By default the API listens on port `8080`, unless changed in `conf/app.conf`.

## Swagger

Generate or regenerate Swagger docs when annotations change:

```bash
swag init
```

Start the server and open the Swagger UI:

```bash
bee run
```

Browse:

```text
http://localhost:8080/swagger/
```

If the path does not load, try:

```text
http://localhost:8080/swagger/index.html
```

## API Reference

### Health Check

```bash
curl http://localhost:8080/api/v1/health
```

### Auth

Register a new user:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret123"}'
```

Login with registered credentials:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'
```

Use the returned `user_id` as the `X-User-ID` header for protected expense routes.

### Expenses

At first, no expense records exist for a new user.
Create the first expense with the returned `user_id`:

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: <user_id>" \
  -d '{"title":"Lunch","amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}'
```

List all expenses for the user:

```bash
curl http://localhost:8080/api/v1/expenses \
  -H "X-User-ID: <user_id>"
```

List expenses with filters:

```bash
curl "http://localhost:8080/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: <user_id>"
```

Sort expenses:

```bash
curl "http://localhost:8080/api/v1/expenses?sort_by=amount&sort_order=desc" \
  -H "X-User-ID: <user_id>"
```

Limit expenses:

```bash
curl "http://localhost:8080/api/v1/expenses?limit=5" \
  -H "X-User-ID: <user_id>"
```

Combined query example:

```bash
curl "http://localhost:8080/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30&sort_by=amount&sort_order=desc&limit=5" \
  -H "X-User-ID: <user_id>"
```

Get a single expense by its ID:

```bash
curl http://localhost:8080/api/v1/expenses/<expense_id> \
  -H "X-User-ID: <user_id>"
```

Update an expense by ID:

```bash
curl -X PUT http://localhost:8080/api/v1/expenses/<expense_id> \
  -H "Content-Type: application/json" \
  -H "X-User-ID: <user_id>" \
  -d '{"title":"Team Dinner","amount":500,"category":"Food","note":"Updated","expense_date":"2025-06-10"}'
```

Delete an expense by ID:

```bash
curl -X DELETE http://localhost:8080/api/v1/expenses/<expense_id> \
  -H "X-User-ID: <user_id>"
```

### Summary

Get expense summary for the user:

```bash
curl http://localhost:8080/api/v1/expenses/summary \
  -H "X-User-ID: <user_id>"
```

Get summary for a date range:

```bash
curl "http://localhost:8080/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: <user_id>"
```

### Allowed categories

The API accepts these categories:

`Food` `Transport` `Housing` `Entertainment` `Shopping` `Healthcare` `Education` `Utilities` `Other`

## Testing

Run the full test suite:

```bash
go test ./...
```

Run tests with coverage enabled:

```bash
go test ./... -cover
```

Generate a detailed coverage report file:

```bash
go test ./... -coverprofile=coverage.out
```

Inspect the coverage summary:

```bash
go tool cover -func=coverage.out
```

Open a browser-based coverage report:

```bash
go tool cover -html=coverage.out
```

If you want coverage across the whole module, include `-coverpkg=./...`:

```bash
go test ./... -coverpkg=./... -coverprofile=coverage.out
```

Current reported coverage: **72.8% of statements** (based on the latest `coverage.out` report generated locally).

### Why these commands

- `go test ./...` runs all package tests in the module.
- `-cover` prints a simple coverage summary.
- `-coverprofile=coverage.out` writes detailed coverage data to a file.
- `go tool cover -func=coverage.out` shows coverage by function.
- `go tool cover -html=coverage.out` opens a browsable HTML coverage report.
- `-coverpkg=./...` measures coverage for all packages in the module.

## Notes

- CSV files are stored in the repository and should exist before running the app.
- Follow `conf/app.conf.example` and provide proper config values in `conf/app.conf`.
- The application resolves CSV locations from config, avoiding hardcoded storage paths.
- If you change Swagger annotations, regenerate docs with `swag init`.
