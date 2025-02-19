# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine

ARG PORT

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build app/main.go

EXPOSE ${PORT}

CMD [ "/app/main" ]