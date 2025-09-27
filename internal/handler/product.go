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

func (ph *productHandler) EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.EditProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process register
	res, err := ph.productService.EditProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.DetailProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process register
	res, err := ph.productService.DetailProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.DeleteProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process register
	res, err := ph.productService.DeleteProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	validationErrors, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.ListProductAdminResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.ListProductAdmin(ctx, req)
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
