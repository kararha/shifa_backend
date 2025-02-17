FROM golang:1.21.13-alpine3.19

WORKDIR /app

# Install dependencies
RUN apk add --no-cache bash mysql-client

# Copy .env file
COPY .env .

# Copy application files
COPY go.mod go.sum ./
COPY . .

RUN go mod download
RUN go build -o main ./cmd/server

EXPOSE 8888

CMD ["./main"]