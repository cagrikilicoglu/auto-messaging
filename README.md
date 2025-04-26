# Auto Messaging System

An automatic message sending system that processes and sends messages from a database at regular intervals.

## Features

- Automatic message sending (processes up to 2 messages every 2 minutes)
  - Message processing starts automatically upon application deployment
  - Processes all unsent messages in the database
  - Messages are processed in chronological order (oldest scheduled messages first)
- Message content character limit validation (500 chars)
- Webhook integration for message delivery
- Database integration for message storage
- Redis caching for message processing
- REST API endpoints for control and monitoring

## Prerequisites

- Go 1.23.0 or higher
- PostgreSQL
- Redis
- Docker and Docker Compose (for containerized deployment)

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
│   ├── database/    # Database utilities
│   └── cache/       # Redis cache implementation
└── config/          # Configuration management
```

## Configuration

The application can be configured using either a YAML file (`config/config.yaml`) or environment variables. The configuration is managed using Viper, which allows for flexible configuration through both methods.

### Configuration File Setup

1. Copy the template configuration files:
```bash
cp config/config.yaml.template config/config.yaml
cp docker-compose.yml.template docker-compose.yml
```

2. Edit `config/config.yaml` with your settings:
```yaml
DB:
  host: localhost
  port: 5432
  user: your_db_user
  password: your_db_password
  name: auto_messaging

redis:
  host: localhost
  port: 6379
  password: your_redis_password
  db: 0

server:
  port: 8080

webhook:
  url: your-webhook-url
  auth_key: your-webhook-auth-key
```

3. Edit `docker-compose.yml` with your settings:
```yaml
# Update the following environment variables in the app service:
- DB_PASSWORD=your_db_password
- REDIS_PASSWORD=your_redis_password
- WEBHOOK_URL=your-webhook-url
- WEBHOOK_AUTH_KEY=your-webhook-auth-key

# Update the following in the postgres service:
- POSTGRES_PASSWORD=your_db_password
```

### Environment Variables

Environment variables can be used to override the YAML configuration. When using Docker, these are set in the `docker-compose.yml` file.

#### Database Configuration
- `DB_HOST`: Database host (default: "postgres" in Docker, "localhost" for local)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: "postgres")
- `DB_PASSWORD`: Database password (required)
- `DB_NAME`: Database name (default: "auto_messaging")

#### Redis Configuration
- `REDIS_HOST`: Redis host (default: "redis" in Docker, "localhost" for local)
- `REDIS_PORT`: Redis port (default: 6379)
- `REDIS_PASSWORD`: Redis password (default: "")
- `REDIS_DB`: Redis database number (default: 0)

#### Server Configuration
- `SERVER_PORT`: API server port (default: 8080)

#### Webhook Configuration
- `WEBHOOK_URL`: URL for the webhook service (required)
- `WEBHOOK_AUTH_KEY`: Authentication key for webhook service (required)

## Installation

1. Clone the repository:
```bash
git clone [repository-url]
cd auto-messaging
```

2. Set up configuration files as described above

3. The application can be run in two ways:

### Using Docker (Recommended)
```bash
docker-compose up -d
```

The Docker environment automatically configures:
- PostgreSQL database
- Redis cache
- Application server
- All necessary environment variables

### Running Locally
1. Make sure PostgreSQL and Redis are running
2. Configure the application through `config/config.yaml` or environment variables
   - For local development, you can keep your `config/config.yaml` with your local settings
   - This file is ignored by git (see `.gitignore`) to prevent committing sensitive information
3. Run the application:
```bash
go run cmd/api/main.go
```

The application will be available at `http://localhost:8080`

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

Note: Message processing starts automatically when the application is deployed. The `/api/v1/messaging/start` endpoint is still available for manual control if needed.

## Message States
- `pending`: Initial state, message waiting to be sent
- `sent`: Message successfully sent
- `failed`: Message sending failed
- `cancelled`: Message was cancelled and won't be sent

## Message Structure
```json
{
  "id": 1,
  "content": "Test message",
  "to": "test@example.com",
  "status": "pending",
  "message_id": "external-message-id",
  "sent_at": "2024-04-26T10:00:00Z",
  "scheduled_at": "2024-04-26T10:00:00Z",
  "created_at": "2024-04-26T09:00:00Z",
  "updated_at": "2024-04-26T09:00:00Z"
}
```

## Webhook Integration
The system sends messages to a configured webhook endpoint. The webhook should expect requests in the following format:

### Request Format
```json
{
  "content": "Test message",
  "to": "test@example.com"
}
```

### Expected Response Format
```json
{
  "message_id": "external-message-id",
  "status": "success"
}
```

## Example Message Creation

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Test message",
    "to": "test@example.com",
    "scheduled_at": "2024-04-26T10:00:00Z"
  }'
```

## Error Handling
The system handles various error scenarios:
- Invalid message content (exceeds 500 characters)
- Invalid email format
- Webhook communication failures
- Database operation failures

Error responses follow this format:
```json
{
  "error": "Error message description"
}
```