FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080 

EXPOSE 1502

CMD ["./main"]
