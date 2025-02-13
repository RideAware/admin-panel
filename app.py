import os
import smtplib
from email.mime.text import MIMEText
from flask import Flask, render_template, request, redirect, url_for, flash, session
from dotenv import load_dotenv
from werkzeug.security import check_password_hash
from database import init_db, get_all_emails, get_admin, create_default_admin

load_dotenv()
app = Flask(__name__)
# Use a secret key from .env; ensure your .env sets SECRET_KEY
app.secret_key = os.getenv('SECRET_KEY')

# SMTP settings (for sending update emails)
SMTP_SERVER = os.getenv('SMTP_SERVER')
SMTP_PORT = int(os.getenv("SMTP_PORT", 465))
SMTP_USER = os.getenv('SMTP_USER')
SMTP_PASSWORD = os.getenv('SMTP_PASSWORD')

# Initialize the database and create default admin user if necessary.
init_db()
create_default_admin()

def login_required(f):
    from functools import wraps
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if "username" not in session:
            return redirect(url_for('login'))
        return f(*args, **kwargs)
    return decorated_function

def process_send_update_email(subject, body):
    """Helper function to send an update email to all subscribers."""
    subscribers = get_all_emails()
    if not subscribers:
        return "No subscribers found."
    try:
        server = smtplib.SMTP_SSL(SMTP_SERVER, SMTP_PORT, timeout=10)
        server.set_debuglevel(True)
        server.login(SMTP_USER, SMTP_PASSWORD)
        for email in subscribers:
            msg = MIMEText(body, 'html', 'utf-8')
            msg['Subject'] = subject
            msg['From'] = SMTP_USER
            msg['To'] = email
            server.sendmail(SMTP_USER, email, msg.as_string())
            print(f"Update email sent to: {email}")
        server.quit()
        return "Email has been sent."
    except Exception as e:
        print(f"Failed to send email: {e}")
        return f"Failed to send email: {e}"

@app.route('/')
@login_required
def index():
    """Displays all subscriber emails"""
    emails = get_all_emails()
    return render_template("admin_index.html", emails=emails)

@app.route('/send_update', methods=['GET', 'POST'])
@login_required
def send_update():
    """Display a form to send an update email; process submission on POST."""
    if request.method == 'POST':
        subject = request.form['subject']
        body = request.form['body']
        # Call the helper function using its new name.
        result_message = process_send_update_email(subject, body)
        flash(result_message)
        return redirect(url_for("send_update"))
    return render_template("send_update.html")

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form.get('username')
        password = request.form.get('password')
        admin = get_admin(username)
        # Expect get_admin() to return a tuple like (username, password_hash)
        if admin and check_password_hash(admin[1], password):
            session['username'] = username
            flash("Logged in successfully", "success")
            return redirect(url_for("index"))
        else:
            flash("Invalid username or password", "danger")
            return redirect(url_for("login"))
    return render_template("login.html")

@app.route('/logout')
def logout():
    session.pop('username', None)
    flash("Logged out successfully", "success")
    return redirect(url_for("login"))

if __name__ == '__main__':
    app.run(port=5000, debug=True)
