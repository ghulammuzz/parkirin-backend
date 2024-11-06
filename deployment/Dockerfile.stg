FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:3.17

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 6012

CMD ["./main"]