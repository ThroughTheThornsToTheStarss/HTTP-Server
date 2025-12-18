from golang:1.24-alpine as builder

workdir /app

copy go.mod go.sum ./
run go mod download

copy . .

run CGO_ENABLED=0 GOOS=linux go build -o http-server ./cmd/server
run CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

from alpine:3.20

workdir /app

copy --from=builder /app/http-server .
copy --from=builder /app/worker .
copy .env .

expose 8080

cmd ["./http-server"]
