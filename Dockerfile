FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

# --- Production Stage ---
FROM alpine:3.20.1 AS prod

# CRITICAL FIX: Install SSL Certs (for NeonDB) and Timezone Data (for KL Time)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=build /app/main ./
COPY --from=build /app/migrations ./migrations

# Set specific port just in case
EXPOSE 8080

CMD ["sh", "-c", "mkdir -p /app/uploads/packages /app/uploads/themes /app/uploads/payment-screenshots && ./main"]