FROM golang:latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/web/*
EXPOSE 4000
ENTRYPOINT ["./main"]
