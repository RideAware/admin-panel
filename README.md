# RideAware Admin Center

This project provides a secure and user-friendly Admin Panel for managing RideAware subscribers and sending update emails. It's designed to work in conjunction with the RideAware landing page application, utilizing a shared database for subscriber management.

## Features


**Secure Admin Authentication:**

* Login protected by username/password authentication using Werkzeug's password hashing.
* Default admin credentials configurable via environment variables.

**Subscriber Management:**
* View a comprehensive list of all subscribed email addresses.

**Email Marketing:**
* Compose and send HTML-rich update emails to all subscribers.
* Supports embedding unsubscribe links in email content for easy opt-out.

**Shared Database:**
* Utilizes a shared PostgreSQL database with the landing page application for consistent subscriber data.

**Centralized Newsletter Storage:**
* Storage of newsletter subject and email bodies in the PostgreSQL database

**Logging:**
* Implemented comprehensive logging throughout the application for better monitoring and debugging.

## Architecture

The Admin Panel is built using Python with the Flask web framework, using the following technologies:

*   **Backend:** Python 3.11+, Flask
*   **Database:** PostgreSQL
*   **Template Engine:** Jinja2
*   **Authentication:** Werkzeug Security
*   **Email:** SMTP (using `smtplib`)
*   **Containerization:** Docker
*   **Configuration:** .env file using `python-dotenv`

## Setup & Deployment

### Prerequisites

*   Docker (recommended for containerized deployment)
*   Python 3.11+ (if running locally without Docker)
*   A PostgreSQL database instance
*   An SMTP account (e.g., SendGrid, Mailgun) for sending emails
*   A `.env` file with configuration details

### .env Configuration

Create a `.env` file in the project root directory with the following environment variables.  Make sure to replace the placeholder values with your actual credentials.

```env
# Flask Application
SECRET_KEY="YourSecretKeyHere" #Used to sign session cookies

# PostgreSQL Database Configuration
PG_HOST=localhost
PG_PORT=5432
PG_DATABASE=rideaware_db
PG_USER=rideaware_user
PG_PASSWORD=rideaware_password

# Admin credentials for the Admin Center
ADMIN_USERNAME=admin
ADMIN_PASSWORD="changeme"  # Change this to a secure password

# SMTP Email Settings
SMTP_SERVER=smtp.example.com
SMTP_PORT=465 #Or another appropriate port
SMTP_USER=your_email@example.com
SMTP_PASSWORD=YourEmailPassword
SENDER_EMAIL=your_email@example.com #Email to send emails from
BASE_URL="your_site_domain.com" # used for unsubscribe links, example.com not https://

#Used for debugging
FLASK_DEBUG=1
```

### Running with Docker (Recommended)

This is the recommended approach for deploying the RideAware Admin Panel

Building the Docker image:
```sh
docker build -t admin-panel .
```

Running the container mapping port 5001:
```sh
docker run -p 5001:5001 admin-panel
```

The application will be accessible at `http://localhost:5001` or `http://<your_server_ip>:5001`

*Note: When running locally with Docker, ensure the .env file is located at the project root. Alternatively, you can pass the variables in the CLI.*

### Running locally

Install Dependencies:

```sh
pip install -r requirements.txt
```

Set Environment Variables:

Ensure all the environment variables specified in the `.env` configuration section are set in your shell environment. You can source the `.env` file:


Then run the admin app:
```sh
python app.py
```

The app will be accessible at `http://127.0.0.1:5001`

### Contributing

Contributions to the RideAware Admin Panel are welcome! Please follow these steps:

* Fork the repository.
* Create a new branch for your feature or bug fix.
* Make your changes and commit them with descriptive commit messages.
* Submit a pull request.