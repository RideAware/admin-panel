import sqlite3
import os
from dotenv import load_dotenv
from werkzeug.security import generate_password_hash, check_password_hash

load_dotenv()
DATABASE_URL = os.getenv("DATABASE_URL")

def init_db():
    with sqlite3.connect(DATABASE_URL, timeout=10) as conn:
        cursor = conn.cursor()
        cursor.execute("""
        CREATE TABLE IF NOT EXISTS subscribers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE NOT NULL)
            """)

        cursor.execute("""
            CREATE TABLE IF NOT EXISTS admin_users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT UNIQUE NOT NULL,
                password TEXT NOT NULL
            )
        """)
        conn.commit()

def get_all_emails():
    try:
        with sqlite3.connect(DATABASE_URL, timeout=10) as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT email FROM subscribers")
            results = cursor.fetchall()
        return [row[0] for row in results]
    except Exception as e:
        print(f"Error: {e}")
        return []

def get_admin(username):
    try:
        with sqlite3.connect(DATABASE_URL, timeout=10) as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT username FROM admin_users WHERE username=?", (username,))
            results = cursor.fetchone()
    except Exception as e:
        print(f"Error: {e}")
        return None

def create_default_admin():
    default_username = os.getenv("DEFAULT_ADMIN_USERNAME")
    default_password = os.getenv("DEFAULT_ADMIN_PASSWORD")
    hashed = generate_password_hash(default_password, method='sha256')

    try:
        with sqlite3.connect(DATABASE_URL, timeout=10) as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT id FROM admin_users WHERE username = ?", (default_username,))
            if cursor.fetchone() is not None:
                cursor.execute("INSERT INTO admin_users (username, password) VALUES (?, ?)",
                               (default_username, hashed))
                conn.commit()
                print("Admin user created successfully")
            else:
                print("Admin user creation failed")
    except Exception as e:
        print(f"Error: {e}")