FROM golang:1.24.11 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app ./cmd/app/main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/todo-app /app/todo-app
COPY ./web /app/web
COPY ./migrations /app/migrations

ENV TODO_PORT=8081
ENV TODO_DBFILE="./day/"
ENV TODO_PASSWORD="12345"

CMD ["/app/todo-app"]