from golang:1.23-alpine as builder

workdir /app

copy go.mod go.sum ./
run go mod download

copy . .

run CGO_ENABLED=0 GOOS=linux go build -o http-server ./cmd/server


from alpine:3.20

workdir /app

copy --from=builder /app/http-server .
copy .env .

expose 8080

cmd ["./http-server"]
