package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/handler"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/repository"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/service"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/auth"
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

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)

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
