# RideAware Admin Center

This project provides an Admin Center for managing the RideAware subscriber list. It connects to the same SQLite database (`subscribers.db`) used by the landing page app (running on port 5000) and allows an administrator to:

- View all currently subscribed email addresses.
- Send overview update emails to all subscribers.
- Unsubscribe emails via the landing page app.

The Admin Center is protected by a login page, and admin credentials (with a salted/hashed password) are stored in the `admin_users` table within the same database.

## Features

- **Admin Login:**  
  Secure login using salted and hashed passwords (via Werkzeug security utilities).

- **Subscriber List:**  
  View all email addresses currently stored in the `subscribers` table.

- **Email Updates:**  
  A form for sending update emails (HTML allowed) to the subscriber list using SMTP.

- **Shared Database:**  
  Both the landing page app (port 5000) and Admin Center (port 5001) connect to the same `subscribers.db`.

## Setup & Running

### Prerequisites

- Docker (for containerized deployment)
- Python 3.11+ (if running locally without Docker)
- An SMTP account (e.g., Spacemail) for sending emails
- A `.env` file with configuration details

### .env Configuration

Create a `.env` file in the project root with the following example variables:

```env
# SMTP settings (shared with the landing page app)
SMTP_SERVER=<email server>
SMTP_PORT=<email port>
SMTP_USER=<email username>
SMTP_PASSWORD=<email password>

# Database file
DATABASE_FILE=subscribers.db

# Admin credentials for the Admin Center
ADMIN_USERNAME=admin
ADMIN_PASSWORD="changeme"  # Change this to a secure password
ADMIN_SECRET_KEY="your_super_secret_key"
```

### Running with Docker

Building the Docker image:
```sh
docker build -t admin-panel .
```

Running the container mapping port 5001:
```sh
docker run -p 5001:5001 admin-panel
```

The app will be accessible at http://ip-address-here:5001

### Running locally
Install the dependencies using **requirements.txt**:

```sh
pip install -r requirements.txt
```

Then run the admin app:
```sh
python app.py
```

The app will be accessible at http://ip-address-here:5001