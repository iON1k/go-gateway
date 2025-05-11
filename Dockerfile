FROM golang:1.22.4 AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download

COPY . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/server


FROM debian:stretch

COPY --from=builder /code/server /usr/local/bin/server
RUN chmod +x /usr/local/bin/server

ENTRYPOINT [ "server" ]