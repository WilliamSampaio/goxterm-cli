FROM golang:1.23-alpine

ARG APP_NAME=goxterm
ARG MAIN_FILE=main.go

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV APP_NAME=${APP_NAME}
ENV MAIN_FILE=${MAIN_FILE}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY main.go .

RUN go build -o bin/${APP_NAME} ${MAIN_FILE}
