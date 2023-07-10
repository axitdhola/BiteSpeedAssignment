FROM golang:1.18.0

WORKDIR /app

COPY . .

RUN go mod tidy
