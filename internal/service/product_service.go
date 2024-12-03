package service

import (
	"context"
	"go-graphql-product-svc/internal/model"
	"go-graphql-product-svc/internal/repository"
	"go-graphql-product-svc/util"
)

type IProductService interface {
	GetAllProduct(ctx context.Context) *[]model.Product
	CreateProduct(ctx context.Context, product model.Product) (*model.Product, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, id string, product model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type ProductService struct {
	Repo repository.IProductRepository
}

// NewProductService creates a new service instance for product-related operations
func NewProductService(repo repository.IProductRepository) *ProductService {
	return &ProductService{
		Repo: repo,
	}
}

func (s *ProductService) GetAllProduct(ctx context.Context) *[]model.Product {
	return s.Repo.Getall(ctx)
}

// CreateProduct calls the repository to create a new product
func (s *ProductService) CreateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	timeNow := util.TimeNow()
	product.CreatedAt = timeNow
	product.UpdatedAt = timeNow
	return s.Repo.Create(ctx, product)
}

// GetProductByID calls the repository to get a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	return s.Repo.FindByID(ctx, id)
}

// UpdateProduct calls the repository to update a product's data
func (s *ProductService) UpdateProduct(ctx context.Context, id string, product model.Product) (*model.Product, error) {
	product.UpdatedAt = util.TimeNow()
	return s.Repo.Update(ctx, id, product)
}

// DeleteProduct calls the repository to delete a product by its ID
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}
