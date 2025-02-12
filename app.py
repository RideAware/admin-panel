import os
import sqlite3
import smtplib
from email.mime.text import MIMEText
from flask import Flask, render_template, request, redirect, url_for, flash
from dotenv import load_dotenv

load_dotenv()
app = Flask(__name__)

app.secret_key = os.getenv('SECRET_KEY')

DATABASE_URL = os.getenv('DATABASE_URL')
SMTP_SERVER = os.getenv('SMTP_SERVER')
SMTP_PORT = int(os.getenv("SMTP_PORT", 465))
SMTP_USER = os.getenv('SMTP_USER')
SMTP_PASSWORD = os.getenv('SMTP_PASSWORD')

def get_all_emails():
    """Retrieve all subscriber emails from the database"""
    try:
        conn = sqlite3.connect(DATABASE_URL)
        cursor = conn.cursor()
        cursor.execute('SELECT email FROM subscribers')
        results = cursor.fetchall()
        conn.close()
        return [row[0] for row in results]
    except Exception as e:
        print(f"Error: {e}")
        return []

def send_update_email(subject, body):
    """Send an update email"""
    subscribers = get_all_emails()
    if not subscribers:
        return "No subscribers found"
    try:
        server = smtplib.SMTP(SMTP_SERVER, SMTP_PORT, timeout=10)
        server.set_debuglevel(True)
        server.login(SMTP_USER, SMTP_PASSWORD)
        for email in subscribers:
            msg = MIMEText(body, 'html', 'utf-8')
            msg['Subject'] = subject
            msg['From'] = SMTP_USER
            msg['To'] = email
            server.sendmail(SMTP_USER, email, msg.as_string())
            print(f"Updated email for {email} has been sent.")
        server.quit()
        return "Email has been sent."
    except Exception as e:
        print(f"Failed to send email: {e}")
        return f"Failed to send email: {e}"

@app.route('/')
def index():
    """Displays all subscriber emails"""
    emails = get_all_emails()
    return render_template("admin_index.html", emails=emails)

@app.route('/send_update_email', methods=['GET', 'POST'])
def send_update_email():
    """Display a form to send an update email"""
    if request.method == 'POST':
        subject = request.form['subject']
        body = request.form['body']
        result_message = send_update_email(subject, body)
        flash(result_message)
        return redirect(url_for("send_update_email"))
    return render_template("send_update.html")

if __name__ == '__main__':
    app.run(port=5001, debug=True)
