## to use with docker compose build this image with `docker build -t mysql-backup .`
FROM golang:1.17.9-alpine3.15

RUN apk update && apk add --no-cache libc6-compat bash mysql-client

RUN mkdir -p /app

COPY ["main.go", "go.mod", "go.sum", ".env", "/app/"]

RUN cd /app && go get -u && go build ./main.go

CMD ["sh", "-c", "cd /app && ./main"]