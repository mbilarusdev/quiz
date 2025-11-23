FROM golang:1.25.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/quiz

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go install github.com/pressly/goose/cmd/goose@latest

RUN ls -l main

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/cmd/quiz/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY .env .
ENV PATH=/usr/local/bin:$PATH
COPY /migrations ./migrations
EXPOSE 8080
CMD ["./main"]