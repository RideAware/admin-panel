FROM python:3.11-slim-buster

# Install build dependencies (build-essential provides gcc and other tools)
RUN apt-get update && apt-get install -y build-essential

WORKDIR /rideaware_landing

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

ENV FLASK_APP=server.py

EXPOSE 5001

CMD ["gunicorn", "--bind", "0.0.0.0:5001", "app:app"]

