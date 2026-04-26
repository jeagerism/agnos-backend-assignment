FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./main.go

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache tzdata ca-certificates
RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /bin/server /app/server

USER app
EXPOSE 8080

ENTRYPOINT ["/app/server"]
