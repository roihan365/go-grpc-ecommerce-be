package service

import (
	"context"
	"os"
	"path/filepath"
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
	DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error)
	EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error)
	DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error)
	ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error)
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

	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("New Product Created"),
		Id: productEntity.Id,
	}, nil
}

func (ps *productService) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// validasi id apakah ada di db
	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	return &product.DetailProductResponse{
		Base: utils.SuccessResponse("Product fetched"),
		Id: productEntity.Id,
		Name: productEntity.Name,
		Description: productEntity.Description,
		Price: productEntity.Price,
		ImageFileName: productEntity.ImageFileName,
	}, nil
}

func (ps *productService) EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error) {
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// validasi id apakah ada di db
	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.EditProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}
	
	// kalau gambar di update, hapus gambar lama
	// if productEntity.ImageFileName != req.ImageFileName {
	// 	newImagePath := filepath.Join("storage", "product", req.ImageFileName)
	// 	_, err := os.Stat(newImagePath)

	// 	if err != nil {
	// 		if os.IsNotExist(err) {
	// 			return &product.EditProductResponse{
	// 				Base: utils.BadRequestResponse("image file not found"),
	// 			}, nil
	// 		}
	// 		return nil, err
	// 	} 

	// 	oldImagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
	// 	err = os.Remove(oldImagePath)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// update ke db
	newProduct := entity.Product{
		Id:           req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:      req.Price,
		ImageFileName: req.ImageFileName,
		UpdatedAt:  time.Now(),
		UpdatedBy:  &claims.FullName,
	}
	err = ps.productRepository.UpdateProduct(ctx, &newProduct)

	if err != nil {
		return nil, err
	}

	return &product.EditProductResponse{
		Base: utils.SuccessResponse("Product has been updated"),
		Id: req.Id,
	}, nil
}

func (ps *productService) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// validasi id apakah ada di db
	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DeleteProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	// update ke db
	err = ps.productRepository.DeleteProduct(ctx, req.Id, time.Now() ,claims.FullName)

	// kalau gambar di update, hapus gambar lama
	if productEntity.ImageFileName != "" && productEntity.ImageFileName != "default.png" {
		oldImagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return &product.DeleteProductResponse{
		Base: utils.SuccessResponse("Product has been deleted"),
	}, nil
}

func (ps *productService) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	products, err := ps.productRepository.ListProductAdmin(ctx)

	if err != nil {
		return nil, err
	}

	var productItems []*product.ListProductAdminResponseItem
	for _, p := range products {
		productItems = append(productItems, &product.ListProductAdminResponseItem{
			Id:            p.Id,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			ImageFileName: p.ImageFileName,
		})
	}

	return &product.ListProductAdminResponse{
		Base: utils.SuccessResponse("Product successfully fetched"),
		Data: productItems,
	}, nil
}

// factory
func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
