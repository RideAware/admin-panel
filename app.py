import os
import logging
import smtplib
from email.mime.text import MIMEText
from flask import (
    Flask,
    render_template,
    request,
    redirect,
    url_for,
    flash,
    session,
)
from dotenv import load_dotenv
from werkzeug.security import check_password_hash
from functools import wraps  # Import wraps
from database import get_connection, init_db, get_all_emails, get_admin, create_default_admin

load_dotenv()
app = Flask(__name__)
app.secret_key = os.getenv("SECRET_KEY")
base_url = os.getenv("BASE_URL")

# SMTP settings (for sending update emails)
SMTP_SERVER = os.getenv("SMTP_SERVER")
SMTP_PORT = int(os.getenv("SMTP_PORT", 465))
SMTP_USER = os.getenv("SMTP_USER")
SMTP_PASSWORD = os.getenv("SMTP_PASSWORD")
SENDER_EMAIL = os.getenv("SENDER_EMAIL", SMTP_USER) # Use SENDER_EMAIL

# Logging setup
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

# Initialize the database and create default admin user if necessary.
init_db()
create_default_admin()

# Decorator for requiring login
def login_required(f):
    @wraps(f)  # Use wraps to preserve function metadata
    def decorated_function(*args, **kwargs):
        if "username" not in session:
            return redirect(url_for("login"))
        return f(*args, **kwargs)

    return decorated_function


def send_update_email(subject, body, email):
    """Sends email, returns True on success, False on failure."""
    try:
        server = smtplib.SMTP_SSL(SMTP_SERVER, SMTP_PORT, timeout=10)
        server.set_debuglevel(False)  # Keep debug level at False for production
        server.login(SMTP_USER, SMTP_PASSWORD)

        unsub_link = f"https://{base_url}/unsubscribe?email={email}"
        custom_body = (
            f"{body}<br><br>"
            f"If you ever wish to unsubscribe, please click <a href='{unsub_link}'>here</a>"
        )

        msg = MIMEText(custom_body, "html", "utf-8")
        msg["Subject"] = subject
        msg["From"] = SENDER_EMAIL  # Use sender email
        msg["To"] = email

        server.sendmail(SENDER_EMAIL, email, msg.as_string())  # Use sender email

        server.quit()
        logger.info(f"Update email sent to: {email}")
        return True
    except Exception as e:
        logger.error(f"Failed to send email to {email}: {e}")
        return False


def process_send_update_email(subject, body):
    """Helper function to send an update email to all subscribers."""
    subscribers = get_all_emails()
    if not subscribers:
        return "No subscribers found."
    try:
        for email in subscribers:
            if not send_update_email(subject, body, email):
                return f"Failed to send to {email}"  # Specific failure message

        # Log newsletter content for audit purposes
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute(
            "INSERT INTO newsletters (subject, body) VALUES (%s, %s)", (subject, body)
        )
        conn.commit()
        cursor.close()
        conn.close()

        return "Email has been sent to all subscribers."
    except Exception as e:
        logger.exception("Error processing sending updates")
        return f"Failed to send email: {e}"


@app.route("/")
@login_required
def index():
    """Displays all subscriber emails"""
    emails = get_all_emails()
    return render_template("admin_index.html", emails=emails)


@app.route("/send_update", methods=["GET", "POST"])
@login_required
def send_update():
    """Display a form to send an update email; process submission on POST."""
    if request.method == "POST":
        subject = request.form["subject"]
        body = request.form["body"]
        result_message = process_send_update_email(subject, body)
        flash(result_message)
        return redirect(url_for("send_update"))
    return render_template("send_update.html")


@app.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form.get("username")
        password = request.form.get("password")
        admin = get_admin(username)
        if admin and check_password_hash(admin[1], password):
            session["username"] = username
            flash("Logged in successfully", "success")
            return redirect(url_for("index"))
        else:
            flash("Invalid username or password", "danger")
            return redirect(url_for("login"))
    return render_template("login.html")


@app.route("/logout")
def logout():
    session.pop("username", None)
    flash("Logged out successfully", "success")
    return redirect(url_for("login"))


if __name__ == "__main__":
    app.run(port=5001, debug=True)
