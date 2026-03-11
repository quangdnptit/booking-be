# Booking API

A simple hotel booking API built with Go and Gin framework.

## Features

- **Rooms Management**: List rooms, get room details, check availability
- **Bookings Management**: Create bookings, list bookings, get booking details
- **Multiple Storage Backends**: In-memory and DynamoDB support
- **Local DynamoDB Development**: Full Docker setup for local testing

## Project Structure

```
booking-be/
├── main.go                 # Application entry point
├── handlers/
│   └── handlers.go        # API endpoint handlers
├── models/
│   └── models.go          # Data models
├── storage/
│   ├── storage.go         # Storage interface and in-memory implementation
│   └── dynamodb.go        # DynamoDB implementation
├── docker-compose.yml     # Docker setup for local DynamoDB
├── .env.example           # Environment variables template
└── README.md              # This file
```

## Quick Start

### 1. In-Memory Storage (Default)

```bash
go run main.go
```

The API will start on `http://localhost:8080`

### 2. DynamoDB with Docker

#### Prerequisites

- Docker and Docker Compose installed
- Go 1.21+

#### Setup Steps

1. **Start local DynamoDB:**

```bash
docker-compose up -d
```

This starts a local DynamoDB instance on `http://localhost:8000`

2. **Create `.env` file:**

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

The `.env` file should contain:
```
STORAGE_TYPE=dynamodb
DYNAMODB_ENDPOINT=http://localhost:8000
PORT=8080
```

3. **Run the application:**

```bash
# Using environment variables
STORAGE_TYPE=dynamodb DYNAMODB_ENDPOINT=http://localhost:8000 go run main.go
```

Or use the `.env` file:
```bash
# For Windows PowerShell
$env:STORAGE_TYPE="dynamodb"
$env:DYNAMODB_ENDPOINT="http://localhost:8000"
go run main.go
```

For bash/Linux:
```bash
source .env
go run main.go
```

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### Rooms
- `GET /rooms` - List all rooms
- `GET /rooms/:id` - Get room details
- `GET /rooms/:id/availability` - Check room availability

### Bookings
- `GET /bookings` - List all bookings
- `POST /bookings` - Create a new booking
- `GET /bookings/:id` - Get booking details

## Example Requests

### Create a Booking

```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "room_id": "room-1",
    "user_id": "user-123",
    "start_date": "2024-01-01",
    "end_date": "2024-01-05"
  }'
```

### List All Rooms

```bash
curl http://localhost:8080/rooms
```

### Get Booking Details

```bash
curl http://localhost:8080/bookings/{booking_id}
```

## Building for Production

```bash
go build -o app main.go
./app
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `STORAGE_TYPE` | `memory` | Storage backend: `memory` or `dynamodb` |
| `DYNAMODB_ENDPOINT` | `` | DynamoDB endpoint URL (for local dev) |
| `PORT` | `8080` | Server port number |
| `AWS_ACCESS_KEY_ID` | `` | AWS credentials (for AWS DynamoDB) |
| `AWS_SECRET_ACCESS_KEY` | `` | AWS credentials (for AWS DynamoDB) |
| `AWS_REGION` | `` | AWS region (for AWS DynamoDB) |

## Development

### Using DynamoDB Locally

To view data in local DynamoDB, install AWS CLI Local:

```bash
npm install -g aws-cli-local dynamodb-local
```

Then query tables:
```bash
awslocal dynamodb scan --table-name rooms
awslocal dynamodb scan --table-name bookings
```

Or use DynamoDB Admin GUI:
```bash
docker run -p 8001:8001 aaronshaf/dynamodb-admin
```
Then open `http://localhost:8001`

### Stop Local DynamoDB

```bash
docker-compose down
```

## Testing

Run tests:
```bash
go test ./...
```

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/google/uuid` - UUID generation
- `github.com/aws/aws-sdk-go-v2` - AWS SDK v2

## License

MIT
