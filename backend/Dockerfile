FROM golang:1.23-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o out/backend ./cmd/backend

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/out/backend .

EXPOSE 8080

CMD ["./backend"]