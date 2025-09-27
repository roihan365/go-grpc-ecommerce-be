package handler

import (
	"context"

	"github.com/roihan365/go-grpc-ecommerce-be/internal/service"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/product"
)

type productHandler struct {
	product.UnimplementedProductServiceServer
	
	productService service.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process register
	res, err := ph.productService.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
}
