package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/grpcmiddleware"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/handler"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/repository"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/service"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/auth"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/product"
	"github.com/roihan365/go-grpc-ecommerce-be/pkg/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	gocache "github.com/patrickmn/go-cache"
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

	cacheService := gocache.New(time.Hour * 24, time.Hour)

	authMiddleware := grpcmiddleware.NewAuthMiddleware(cacheService)

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware, authMiddleware.Middleware),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductServiceServer(serv, productHandler)

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
