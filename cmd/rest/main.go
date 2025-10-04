package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/handler"
)

func main() {
	app := fiber.New()

	app.Post("product/upload-image", handler.UploadProductImage)

	app.Listen(":3000")
}