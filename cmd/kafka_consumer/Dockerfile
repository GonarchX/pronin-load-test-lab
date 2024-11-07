FROM golang:1.23-alpine as builder

WORKDIR /build

COPY ../../go.mod go.mod

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /server ./cmd/kafka_consumer/kafka_consumer.go

FROM scratch
COPY --from=builder server /bin/server

ENTRYPOINT ["/bin/server"]