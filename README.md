# Auto Messaging System

An automatic message sending system that processes and sends messages from a database at regular intervals.

## Features

- Automatic message sending (2 messages every 2 minutes)
- Message content character limit validation (500 chars)
- Webhook integration for message delivery
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
│   ├── client/      # External service clients (webhook)
│   ├── controller/  # Business logic
│   ├── handler/     # HTTP handlers
│   ├── model/       # Data models
│   ├── repository/  # Database operations
│   └── router/      # API route definitions
├── pkg/
│   └── database/    # Database utilities
└── config/          # Configuration management
```

## Configuration

The application is configured using environment variables:

### Database Configuration
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: auto_messaging)

### Server Configuration
- `API_PORT`: API server port (default: 8080)

### Webhook Configuration
- `WEBHOOK_URL`: URL for the webhook service (required)
- `WEBHOOK_AUTH_KEY`: Authentication key for webhook service (required)

## Installation

1. Clone the repository:
```bash
git clone [repository-url]
cd auto-messaging
```

2. Create a `.env` file with your configuration (see .env.example)

3. Start the services using Docker Compose:
```bash
docker-compose up -d
```

4. The application will be available at `http://localhost:8080`

## API Documentation

Once the application is running, you can access the Swagger documentation at:
`http://localhost:8080/swagger/index.html`

## API Endpoints

### Message Management
- `POST /api/v1/messages` - Create a new message
- `GET /api/v1/messages` - Get all messages
- `GET /api/v1/messages/{id}` - Get a specific message
- `PUT /api/v1/messages/{id}/status` - Update message status

### Message Processing Control
- `POST /api/v1/messaging/start` - Start automatic message sending
- `POST /api/v1/messaging/stop` - Stop automatic message sending
- `GET /api/v1/messaging/sent` - Get list of sent messages

## Message States
- `pending`: Initial state, message waiting to be sent
- `sent`: Message successfully sent
- `failed`: Message sending failed