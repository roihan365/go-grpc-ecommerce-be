package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/handler"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/service"
	"github.com/roihan365/go-grpc-ecommerce-be/pkg/database"
	"github.com/roihan365/go-grpc-ecommerce-be/pkg/grpcmiddleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	godotenv.Load(".env")
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error when listening: %v", err)
	}

	database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database")

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware),
	)

	service.RegisterHelloWorldServer(serv, serviceHandler)

	// tidak disarankan untuk digunakan di production
	if os.Getenv("ENVIRONTMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection registered")
	}

	log.Println("Server is running on port 50052")

	if err := serv.Serve(lis); err != nil {
		log.Panicf("Server is error %v", err)
	}
}
