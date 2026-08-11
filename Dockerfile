FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the server executable (CGO_ENABLED=0 ensures it runs smoothly on Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o run_chatbot ./server

FROM alpine:latest

WORKDIR /root/

# Copy the compiled server from Stage 1
COPY --from=builder /app/run_chatbot .

EXPOSE 8080

CMD ["./run_chatbot"]