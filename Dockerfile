FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o break-tester main.go
RUN CGO_ENABLED=0 go build -o break-tester main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/break-tester .
COPY config.json .

ENTRYPOINT ["./break-tester"]