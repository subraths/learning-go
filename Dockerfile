# syntax=docker/dockerfile:experimental
FROM golang:alpine
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
