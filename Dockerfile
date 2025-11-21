# Stage 1: Build
FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo \
    -o admin-panel ./cmd/admin-panel

# Stage 2: Runtime
FROM docker.io/library/alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/admin-panel .
COPY .env .env
COPY web ./web

EXPOSE 5001

CMD ["./admin-panel"]
