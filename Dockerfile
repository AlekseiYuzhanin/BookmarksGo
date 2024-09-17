FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o bookmarkSearch ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bookmarkSearch .

RUN test -f ./bookmarkSearch && chmod +x ./bookmarkSearch

# Открываем порт
EXPOSE 8080

CMD ["./bookmarkSearch"]