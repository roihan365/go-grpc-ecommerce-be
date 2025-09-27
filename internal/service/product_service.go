package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/entity"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/repository"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/product"

	jwtEntity "github.com/roihan365/go-grpc-ecommerce-be/internal/entity/jwt"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func (ps *productService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	
	// cek dulu apakah user admin
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// cek apakah ada image product

	// insert ke db


	// response

	// cek email ke db
	productEntity := entity.Product{
		Id:           uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Price:      req.Price,
		ImageFileName: req.ImageFileName,
		CreatedAt:  time.Now(),
		CreatedBy:  claims.FullName,
	}
	
	err = ps.productRepository.CreateProduct(ctx, &productEntity)


	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("New Product Created"),
		Id: productEntity.Id,
	}, nil
}
// factory
func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
