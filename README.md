# Auto Messaging System

An automatic message sending system that processes and sends messages from a database at regular intervals.

## Features

- Automatic message sending every 2 minutes
- Database integration for message storage
- REST API endpoints for control and monitoring

## Prerequisites

- Go 1.21 or higher
- PostgreSQL

## Project Structure

```
auto-messaging/
├── cmd/
│   └── api/         # API server implementation
├── internal/
│   ├── handler/     # HTTP handlers
│   ├── model/       # Data models
│   ├── repository/  # Database operations
│   └── service/     # Business logic
├── pkg/
│   └── database/    # Database utilities
└── config/          # Configuration management
```

## Configuration

The application can be configured using environment variables:

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `API_PORT`: API server port (default: 8080)

## Installation

1. Clone the repository:
```bash
git clone [repository-url]
cd auto-messaging
```

2. Start the services using Docker Compose:
```bash
docker-compose up -d
```

3. The application will be available at `http://localhost:8080`

## API Documentation

Once the application is running, you can access the Swagger documentation at:
`http://localhost:8080/swagger/index.html`

## API Endpoints

- `POST /api/v1/messaging/start` - Start automatic message sending
- `POST /api/v1/messaging/stop` - Stop automatic message sending
- `GET /api/v1/messaging/sent` - Get list of sent messages