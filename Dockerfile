FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.work ./

WORKDIR /app/avs
COPY ./avs/go.mod ./avs/go.sum ./

WORKDIR /app/shared
COPY ./shared/go.mod ./shared/go.sum ./

WORKDIR /app/mock_secware
COPY ./mock_secware/go.mod ./mock_secware/go.sum ./

WORKDIR /app
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /app/avs/avs goplus/avs/cmd

FROM docker:latest

ARG DOCKER_USER
ARG DOCKER_PWD

RUN docker login -u $DOCKER_USER -p $DOCKER_PWD

COPY --from=builder /app/avs/avs /app/avs
RUN chmod +x /app/avs

CMD ["/bin/sh", "-c", "/app/avs"]

