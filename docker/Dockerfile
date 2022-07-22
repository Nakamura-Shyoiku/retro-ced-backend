FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o build/main cmd/server/main.go 
ENV RUN_ENV=docker_dev

EXPOSE 8080

CMD [ "build/main" ]