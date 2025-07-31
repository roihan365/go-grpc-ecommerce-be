package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/handler"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	godotenv.Load(".env")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error when listening: %v", err)
	}

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer()

	service.RegisterHelloWorldServer(serv, serviceHandler)

	// tidak disarankan untuk digunakan di production
	if os.Getenv("ENVIRONTMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection registered")
	}

	log.Println("Server is running on port 50051")

	if err := serv.Serve(lis); err != nil {
		log.Panicf("Server is error %v", err)
	}
}
