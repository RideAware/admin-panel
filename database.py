import os
import logging
import psycopg2
from psycopg2 import IntegrityError
from dotenv import load_dotenv
from werkzeug.security import generate_password_hash

load_dotenv()

# Logging setup
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


def get_connection():
    """Return a new connection to the PostgreSQL database."""
    try:
        conn = psycopg2.connect(
            host=os.getenv("PG_HOST"),
            port=os.getenv("PG_PORT"),
            dbname=os.getenv("PG_DATABASE"),
            user=os.getenv("PG_USER"),
            password=os.getenv("PG_PASSWORD"),
            connect_timeout=10,
        )
        return conn
    except Exception as e:
        logger.error(f"Database connection error: {e}")
        raise


def init_db():
    """Initialize the database tables."""
    conn = None
    try:
        conn = get_connection()
        cursor = conn.cursor()

        # Create subscribers table (if not exists)
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS subscribers (
                id SERIAL PRIMARY KEY,
                email TEXT UNIQUE NOT NULL
            )
        """
        )

        # Create admin_users table (if not exists)
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS admin_users (
                id SERIAL PRIMARY KEY,
                username TEXT UNIQUE NOT NULL,
                password TEXT NOT NULL
            )
        """
        )

        # Newsletter storage
        cursor.execute(
            """
        CREATE TABLE IF NOT EXISTS newsletters (
            id SERIAL PRIMARY KEY,
            subject TEXT NOT NULL,
            body TEXT NOT NULL,
            sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
        """
        )

        conn.commit()
        logger.info("Database initialized successfully.")
    except Exception as e:
        logger.error(f"Database initialization error: {e}")
        if conn:
            conn.rollback()  # Rollback if there's an error

        raise
    finally:
        if conn:
            cursor.close()
            conn.close()


def get_all_emails():
    """Return a list of all subscriber emails."""
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("SELECT email FROM subscribers")
        results = cursor.fetchall()
        emails = [row[0] for row in results]
        logger.debug(f"Retrieved emails: {emails}")
        return emails
    except Exception as e:
        logger.error(f"Error retrieving emails: {e}")
        return []
    finally:
        if conn:
            cursor.close()
            conn.close()


def add_email(email):
    """Insert an email into the subscribers table."""
    conn = None
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("INSERT INTO subscribers (email) VALUES (%s)", (email,))
        conn.commit()
        logger.info(f"Email {email} added successfully.")
        return True
    except IntegrityError:
        logger.warning(f"Attempted to add duplicate email: {email}")
        return False
    except Exception as e:
        logger.error(f"Error adding email {email}: {e}")
        return False
    finally:
        if conn:
            cursor.close()
            conn.close()


def remove_email(email):
    """Remove an email from the subscribers table."""
    conn = None
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("DELETE FROM subscribers WHERE email = %s", (email,))
        rowcount = cursor.rowcount
        conn.commit()
        logger.info(f"Email {email} removed successfully.")
        return rowcount > 0
    except Exception as e:
        logger.error(f"Error removing email {email}: {e}")
        return False
    finally:
        if conn:
            cursor.close()
            conn.close()


def get_admin(username):
    """Retrieve admin credentials for a given username.
    Returns a tuple (username, password_hash) if found, otherwise None.
    """
    conn = None
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute(
            "SELECT username, password FROM admin_users WHERE username = %s",
            (username,),
        )
        result = cursor.fetchone()
        return result  # (username, password_hash)
    except Exception as e:
        logger.error(f"Error retrieving admin: {e}")
        return None
    finally:
        if conn:
            cursor.close()
            conn.close()


def create_default_admin():
    """Create a default admin user if one doesn't already exist."""
    default_username = os.getenv("ADMIN_USERNAME", "admin")
    default_password = os.getenv("ADMIN_PASSWORD", "changeme")
    hashed_password = generate_password_hash(default_password, method="pbkdf2:sha256")
    conn = None
    try:
        conn = get_connection()
        cursor = conn.cursor()

        # Check if the admin already exists
        cursor.execute(
            "SELECT id FROM admin_users WHERE username = %s", (default_username,)
        )
        if cursor.fetchone() is None:
            cursor.execute(
                "INSERT INTO admin_users (username, password) VALUES (%s, %s)",
                (default_username, hashed_password),
            )
            conn.commit()
            logger.info("Default admin created successfully")
        else:
            logger.info("Default admin already exists")
    except Exception as e:
        logger.error(f"Error creating default admin: {e}")
        if conn:
            conn.rollback()
    finally:
        if conn:
            cursor.close()
            conn.close()

