FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
# Don't copy .env in production
# COPY .env .

# Add this for debugging
RUN apk add --no-cache curl

EXPOSE 8080

# Explicitly set the PORT environment variable
#ENV PORT=8080

# Add a healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/health || exit 1

CMD ["./main"]
