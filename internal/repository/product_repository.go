package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/roihan365/go-grpc-ecommerce-be/internal/entity"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductById(ctx context.Context, id string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
}

type productRepository struct {
	db sql.DB
}

func (pr *productRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	_, err := pr.db.ExecContext(
		ctx,
		"INSERT INTO products (id, name, description, price, image_file_name, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		product.Id,
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.CreatedAt,
		product.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	var productEntity entity.Product
	row := pr.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, price, image_file_name FROM products WHERE id=$1 AND is_deleted=false",
		id,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&productEntity.Id,
		&productEntity.Name,
		&productEntity.Description,
		&productEntity.Price,
		&productEntity.ImageFileName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &productEntity, nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	_, err := pr.db.ExecContext(
		ctx,
		"UPDATE products SET name=$1, description=$2, price=$3, image_file_name=$4, updated_at=$5, updated_by=$6 WHERE id=$7",
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.UpdatedAt,
		product.UpdatedBy,
		product.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error {
	_, err := pr.db.ExecContext(
		ctx,
		"UPDATE products SET is_deleted=true, deleted_at=$1, deleted_by=$2 WHERE id=$3",
		deletedAt,
		deletedBy,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{
		db: *db,
	}
}
