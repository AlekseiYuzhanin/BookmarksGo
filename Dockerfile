FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bookmarkSearch ./cmd/main.go

RUN chmod +x ./bookmarkSearch

EXPOSE 8080

CMD ["./bookmarkSearch"]
