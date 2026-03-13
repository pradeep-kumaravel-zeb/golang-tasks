# Student Course Enrollment System

A Go-based REST API for managing student course enrollments with automatic payment status calculation.

## Database Configuration

- Host: localhost
- Port: 5432
- User: postgres
- Password: root
- Database: PostgreSQL 16
- Schema: public
- SSL Mode: disable

## Setup

1. Install dependencies:
```bash
cd development
go mod download
```

2. Run the application:
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Enrollments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/enrollments` | Get all enrollments (supports filters: `course_name`, `payment_status`, `student_name`) |
| GET | `/enrollments/:id` | Get enrollment by ID |
| POST | `/enrollments` | Create new enrollment |
| PUT | `/enrollments/:id` | Fully update enrollment |
| PATCH | `/enrollments/:id` | Partially update enrollment |
| DELETE | `/enrollments/:id` | Delete enrollment (only if unpaid or course not started) |

## Features

- Auto-migration: Tables are created automatically on startup
- Payment status auto-calculation based on fee paid vs course fee
- Email validation
- Student and course management (auto-created when needed)
- Query filtering support

## Project Structure

```
development/
├── cmd/
│   └── main.go                 # Application entry point
├── pkg/
│   ├── server/
│   │   └── server.go           # Server initialization and routing
│   └── database/
│       └── database.go         # Database connection and migration
├── internal/
│   ├── config/
│   │   └── config.go           # Database URL builder
│   ├── models/
│   │   └── model.go            # GORM models
│   ├── dto/
│   │   └── enrollment_dto.go   # Request/Response DTOs
│   ├── repository/
│   │   └── repository.go       # Database operations
│   ├── service/
│   │   └── service.go          # Business logic
│   ├── handlers/
│   │   └── handler.go          # HTTP handlers
│   └── utils/
│       └── utils.go            # Validation utilities
└── go.mod
```

## Note

Before creating enrollments, you need to create courses first. The system expects courses to exist in the database.
