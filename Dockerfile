FROM golang:1.19-alpine AS builder
WORKDIR /build

COPY . .
RUN go build -o ./output/bot ./cmd/bot/main.go

FROM golang:1.19-alpine AS runner
WORKDIR /app

COPY --from=builder /build/output/bot .

CMD ["./bot"]