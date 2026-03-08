# Todo List API (Golang)

A simple, robust Todo List API built with Go, featuring a central schema, SQLite database, and automated Swagger documentation.

## Features

- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/) with SQLite
- **Documentation**: [Swaggo](https://github.com/swaggo/swag) for automated Swagger UI and `swagger.json` generation
- **Automation**: Makefile for easy building and running

## Prerequisites

- [Go](https://golang.org/doc/install) (1.18+)
- `make` (optional, but recommended)

## Getting Started

### 1. Install Dependencies

```bash
go mod download
```

### 2. Run the Application

This will automatically generate/update the Swagger documentation and start the server:

```bash
make run
```

The server will start at `http://localhost:8080`.

### 3. Build the Application

To generate documentation and build the binary:

```bash
make build
```

## API Documentation

The API includes built-in Swagger UI and a health check endpoint for exploration and testing.

- **Health Check**: [http://localhost:8080/ping](http://localhost:8080/ping) (Returns `{"message": "pong"}`)
- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Swagger JSON**: `docs/swagger.json`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/todos` | List all todos |
| POST | `/api/v1/todos` | Create a new todo |
| GET | `/api/v1/todos/:id` | Get a specific todo |
| PUT | `/api/v1/todos/:id` | Update a todo |
| DELETE | `/api/v1/todos/:id` | Delete a todo |

## Project Structure

- `main.go`: API logic, routing, and database initialization.
- `models/todo.go`: Central schema for the Todo item.
- `docs/`: Generated Swagger documentation.
- `Makefile`: Commands for building and running the project.
- `todos.db`: SQLite database file (created automatically).

## Updating Schema & Swagger

Whenever you modify the `models/todo.go` or update the API annotations in `main.go`, simply run:

```bash
make swag-init
```

This will refresh `docs/swagger.json` and `docs/docs.go`.
