FROM golang:1.24.1-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
RUN go get github.com/pressly/goose/cmd/goose
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder  /app/migrations ./migrations
COPY --from=builder /app/scripts/wait_for_db.sh ./scripts/wait_for_db.sh
RUN chmod +x ./scripts/wait_for_db.sh 
EXPOSE 8080
CMD ["./main"]
