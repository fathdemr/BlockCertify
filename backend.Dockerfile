# ── Build Stage ──────────────────────────────────────────────────────────────
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ── Final Stage ───────────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/server .

# Arweave wallet key (needed for on-chain transactions)
COPY arweave_keyfile.json .

# Static public assets
COPY public/ ./public/

# Uploads directory (mount as volume for persistence)
RUN mkdir -p /app/uploads
VOLUME /app/uploads

EXPOSE 8080

CMD ["./server"]

