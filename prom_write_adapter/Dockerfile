FROM golang:1.14-alpine

ENV GO111MODULE=on

WORKDIR /app/server

COPY go.mod go.sum  ./

RUN go mod download
COPY . .

RUN go build

CMD ["./prom_write_adapter"]
