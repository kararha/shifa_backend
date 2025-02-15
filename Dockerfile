FROM golang:1.21.13-alpine3.19

WORKDIR /app

# Add wait-for-it script
RUN apk add --no-cache bash mysql-client

COPY go.mod go.sum ./
COPY . .

RUN go mod download
RUN go build -o main ./cmd/server

EXPOSE 8888

# Add healthcheck
HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=3 \
    CMD mysqladmin ping -h mysql -u${DB_USER} -p${DB_PASSWORD} || exit 1

CMD ["./main"]