FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o run_chatbot ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/run_chatbot .

EXPOSE 8080

CMD ["./run_chatbot"]