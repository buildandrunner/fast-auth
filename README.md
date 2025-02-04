# Space Auth

## Overview
Space Auth is a lightweight authentication microservice written in Go, utilizing Gin as the web framework and Redis as the session store. It provides endpoints for user registration, login, logout, and session management.

## Features
- User registration with Argon2id password hashing.
- Secure login with password validation.
- Session management using Redis.
- REST API with JSON responses.
- Dockerized for deployment.

## Requirements
- Go 1.23 or later
- Redis server
- Docker (optional, for containerized deployment)

## Installation

### Clone the repository:
```sh
git clone git@github.com:mar-cial/space-auth.git
cd space-auth
```

### Install dependencies:
```sh
go mod tidy
```

## Configuration
Set the required environment variable:
```sh
export REDIS_URL=redis://localhost:6379
```

## Running the Service

### Locally
```sh
go run cmd/main.go
```

### Using Docker
```sh
docker build -t space-auth .
docker run -p 8080:8080 -e REDIS_URL=redis://your-redis-host:6379 space-auth
```

## API Endpoints

### Register a User
```
POST /register
{
  "phonenumber": "+1234567890",
  "password": "securepassword"
}
```

### Login
```
POST /login
{
  "phonenumber": "+1234567890",
  "password": "securepassword"
}
```

### Logout
```
POST /logout
{
  "token": "session-token"
}
```

## Project Structure
```
space-auth/
├── cmd/main.go                # Entry point
├── internal/
│   ├── adapter/
│   │   ├── handler/           # HTTP handlers
│   │   ├── repository/redis/  # Redis repository
│   ├── core/
│   │   ├── domain/            # Domain entities
│   │   ├── port/              # Interfaces
│   │   ├── service/           # Business logic
├── templates/                 # HTML templates
├── Dockerfile                 # Docker configuration
├── go.mod                     # Go module file
```

## License
This project is licensed under the MIT License.

