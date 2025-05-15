FROM golang:1.24.1-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/migrations ./migrations
COPY scripts/wait_for_db.sh .

RUN echo $'#!/bin/sh\n\
/app/wait_for_db.sh\n\
echo "Running migrations..."\n\
goose -dir ./migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" up\n\
echo "Starting application with Swagger docs..."\n\
exec /app/main' > start.sh && chmod +x start.sh

RUN chmod +x /app/wait_for_db.sh && \
    chmod +x /app/start.sh

EXPOSE 8080
CMD ["./start.sh"]