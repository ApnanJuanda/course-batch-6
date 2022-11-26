# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build app/main.go

EXPOSE 1234

CMD [ "/app/main" ]