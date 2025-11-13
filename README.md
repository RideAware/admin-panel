# RideAware Admin Center

This project provides a secure and user-friendly Admin Panel for managing RideAware subscribers and sending update emails. It's designed to work in conjunction with the RideAware landing page application, utilizing a shared database for subscriber management.

## Features

**Secure Admin Authentication:**
* Login protected by username/password authentication using bcrypt password hashing.
* Default admin credentials configurable via environment variables.
* Session-based authentication with secure HTTP-only cookies.

**Subscriber Management:**
* View a comprehensive list of all subscribed email addresses.

**Email Marketing:**
* Compose and send HTML-rich update emails to all subscribers.
* Supports embedding unsubscribe links in email content for easy opt-out.

**Shared Database:**
* Utilizes a shared PostgreSQL database with the landing page application for consistent subscriber data.

**Centralized Newsletter Storage:**
* Storage of newsletter subject and email bodies in the PostgreSQL database.

**Comprehensive Logging:**
* Implemented structured logging throughout the application for better monitoring and debugging.

## Architecture

The Admin Panel is built using Go with the Gin web framework, using the following technologies:

* **Backend:** Go 1.23+, Gin Web Framework
* **Database:** PostgreSQL with `lib/pq` driver
* **Authentication:** Bcrypt for password hashing, Gorilla Sessions for session management
* **Email:** SMTP via `go-mail` library
* **Containerization:** Podman/Docker with multi-stage builds
* **Configuration:** Environment variables via `godotenv`

## Setup & Deployment

### Prerequisites

* Podman or Docker (recommended for containerized deployment)
* Go 1.23+ (if running locally without containers)
* A PostgreSQL database instance
* An SMTP account (e.g., SendGrid, Gmail, Mailgun) for sending emails
* A `.env` file with configuration details

### .env Configuration

Create a `.env` file in the project root directory with the following environment variables. Make sure to replace the placeholder values with your actual credentials.

```env
# Go Application
PORT=5001
GIN_MODE=debug  # Use 'release' in production

# PostgreSQL Database Configuration
PG_HOST=localhost
PG_PORT=5432
PG_DATABASE=rideaware
PG_USER=postgres
PG_PASSWORD=your_postgres_password

# Admin credentials for the Admin Center
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_password  # Change this to a secure password
SECRET_KEY=your_secret_key_here      # Used to sign session cookies

# SMTP Email Settings
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=465                        # Or another appropriate port
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password      # Use app-specific password for Gmail
SENDER_EMAIL=your_email@gmail.com    # Email address to send from

# Application Settings
BASE_URL=example.com                 # Used for unsubscribe links (without https://)
```

### Running with Podman (Recommended)

This is the recommended approach for deploying the RideAware Admin Panel.

**Building the image:**
```sh
podman build -t admin-panel:latest .
```

**Running the container:**
```sh
podman run -d \
  --name admin-panel \
  -p 5001:5001 \
  --env-file .env \
  admin-panel:latest
```

**Viewing logs:**
```sh
podman logs -f admin-panel
```

The application will be accessible at `http://localhost:5001` or `http://<your_server_ip>:5001`

### Running with Podman Compose

Create a `podman-compose.yml`:

```yaml
version: '3.8'

services:
  admin-panel:
    build: .
    ports:
      - "5001:5001"
    env_file:
      - .env
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: docker.io/library/postgres:15-alpine
    environment:
      POSTGRES_USER: ${PG_USER}
      POSTGRES_PASSWORD: ${PG_PASSWORD}
      POSTGRES_DB: ${PG_DATABASE}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

Then run:
```sh
podman-compose up -d
podman-compose logs -f admin-panel
```

### Running Locally (Development)

**Install Go 1.23+**

Ensure you have Go installed. Download from [golang.org](https://golang.org/dl)

**Install dependencies:**
```sh
go mod tidy
```

**Run the application:**
```sh
go run ./cmd/admin-panel
```

The app will be accessible at `http://127.0.0.1:5001`

**For hot-reloading during development, install `air`:**
```sh
go install github.com/cosmtrek/air@latest
air
```

## API Endpoints

| Method | Endpoint | Description | Protected |
|--------|----------|-------------|-----------|
| GET | `/login` | Login page | No |
| POST | `/login` | Submit login credentials | No |
| GET | `/logout` | Logout and clear session | Yes |
| GET | `/` | View subscriber list | Yes |
| GET | `/send_update` | Newsletter compose form | Yes |
| POST | `/send_update` | Send newsletter to all subscribers | Yes |

## Development

### Adding New Features

1. Create a new handler in `internal/handlers/`
2. Add routes in `cmd/admin-panel/main.go`
3. Add any new dependencies to `go.mod` via `go get`
4. Run `go mod tidy` to sync dependencies

### Database Migrations

The application automatically creates required tables on startup if they don't exist. To add new tables, modify the `createTables()` function in `internal/database/database.go`.

### Environment Variables

Configuration is centralized in `internal/config/config.go`. All environment variables are loaded via the `godotenv` package. Defaults are provided for development.

## Contributing

Contributions to the RideAware Admin Panel are welcome! Please follow these steps:

* Fork the repository.
* Create a new branch for your feature or bug fix.
* Make your changes and commit them with descriptive commit messages.
* Ensure code builds with `go build ./cmd/admin-panel`
* Run `go fmt ./...` to format code
* Submit a pull request.

## Troubleshooting

**Database Connection Failed**
* Verify PostgreSQL is running and accessible
* Check `PG_HOST`, `PG_PORT`, `PG_USER`, and `PG_PASSWORD` in `.env`
* Test connection manually: `psql -h <host> -U <user> -d <database>`

**Email Not Sending**
* Verify SMTP credentials are correct
* Check `SMTP_SERVER` and `SMTP_PORT` match your email provider
* For Gmail, use an [app-specific password](https://support.google.com/accounts/answer/185833)
* Check container logs for SMTP errors

**Session/Login Issues**
* Ensure `SECRET_KEY` is set and consistent
* Clear browser cookies and try again
* Check that `GIN_MODE` is not set to `release` in development