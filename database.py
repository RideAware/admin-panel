import os
import psycopg2
from psycopg2 import IntegrityError
from dotenv import load_dotenv
from werkzeug.security import generate_password_hash
load_dotenv()

def get_connection():
    """Return a new connection to the PostgreSQL database."""
    return psycopg2.connect(
        host=os.getenv("PG_HOST"),
        port=os.getenv("PG_PORT"),
        dbname=os.getenv("PG_DATABASE"),
        user=os.getenv("PG_USER"),
        password=os.getenv("PG_PASSWORD"),
        connect_timeout=10
    )

def init_db():
    """Initialize the database tables."""
    conn = get_connection()
    cursor = conn.cursor()
    # Create subscribers table (if not exists)
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS subscribers (
            id SERIAL PRIMARY KEY,
            email TEXT UNIQUE NOT NULL
        )
    """)
    # Create admin_users table (if not exists)
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS admin_users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        )
    """)
    conn.commit()
    cursor.close()
    conn.close()

def get_all_emails():
    """Return a list of all subscriber emails."""
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("SELECT email FROM subscribers")
        results = cursor.fetchall()
        cursor.close()
        conn.close()
        return [row[0] for row in results]
    except Exception as e:
        print(f"Error retrieving emails: {e}")
        return []

def add_email(email):
    """Insert an email into the subscribers table."""
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("INSERT INTO subscribers (email) VALUES (%s)", (email,))
        conn.commit()
        cursor.close()
        conn.close()
        return True
    except IntegrityError:
        return False
    except Exception as e:
        print(f"Error adding email: {e}")
        return False

def remove_email(email):
    """Remove an email from the subscribers table."""
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("DELETE FROM subscribers WHERE email = %s", (email,))
        conn.commit()
        rowcount = cursor.rowcount
        cursor.close()
        conn.close()
        return rowcount > 0
    except Exception as e:
        print(f"Error removing email: {e}")
        return False

def get_admin(username):
    """Retrieve admin credentials for a given username.
       Returns a tuple (username, password_hash) if found, otherwise None.
    """
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("SELECT username, password FROM admin_users WHERE username = %s", (username,))
        result = cursor.fetchone()
        cursor.close()
        conn.close()
        return result  # (username, password_hash)
    except Exception as e:
        print(f"Error retrieving admin: {e}")
        return None

def create_default_admin():
    """Create a default admin user if one doesn't already exist."""
    default_username = os.getenv("ADMIN_USERNAME", "admin")
    default_password = os.getenv("ADMIN_PASSWORD", "changeme")
    hashed = generate_password_hash(default_password, method="pbkdf2:sha256")
    try:
        conn = get_connection()
        cursor = conn.cursor()
        # Check if the admin already exists
        cursor.execute("SELECT id FROM admin_users WHERE username = %s", (default_username,))
        if cursor.fetchone() is None:
            cursor.execute("INSERT INTO admin_users (username, password) VALUES (%s, %s)",
                           (default_username, hashed))
            conn.commit()
            print("Default admin created successfully")
        else:
            print("Default admin already exists")
        cursor.close()
        conn.close()
    except Exception as e:
        print(f"Error creating default admin: {e}")
