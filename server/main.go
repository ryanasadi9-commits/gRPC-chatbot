package main

import (
	"database/sql"
	"log"
	"net"

	"hamrahTask1/proto"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// connecting to the database and redis
	db, err := sql.Open("pgx", "postgres://ryanasadi:12345678@db:5432/chatbot_db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to PostgreSQL.")

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer rdb.Close()
	log.Println("Successfully connected to Redis.")

	// waiting to make the tcp connection
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterChatbotServiceServer(grpcServer, &server{db: db, rdb: rdb})

	log.Printf("gRPC server is successfully running and listening on port :8080...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
