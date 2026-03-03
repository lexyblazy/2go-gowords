FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/server.go 

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config.json .
COPY --from=builder /app/dictionary.txt .


EXPOSE 8080
CMD ["./server"]