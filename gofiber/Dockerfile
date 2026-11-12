FROM golang:1.23-alpine AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gofiber-template ./cmd/api

# ใช้ alpine แทน debian
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget
ENV TZ=Asia/Bangkok

COPY --from=builder /app/gofiber-template .
COPY --from=builder /app/pkg/config ./pkg/config

# ไม่ต้อง copy .env ใน production - ใช้ environment variables
# COPY .env .env

RUN mkdir -p /app/uploads

EXPOSE 9095

CMD ["./gofiber-template"]