# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOTOOLCHAIN=auto
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Run stage
FROM alpine:3.19
WORKDIR /app
ENV PORT=8080
# DB envs should be provided via compose/env
COPY --from=builder /app/server /app/server
COPY --from=builder /app/internal /app/internal
COPY --from=builder /app/migrations /app/migrations
EXPOSE 8080
CMD ["/app/server"]
