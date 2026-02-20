# build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

# runtime stage
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/app .

ENV PORT=8080
EXPOSE $PORT

CMD ["./app"]