package repository

import (
	"context"
	"fmt"
	"go-graphql-product-svc/internal/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IProductRepository interface {
	Getall(ctx context.Context) *[]model.Product
	Create(ctx context.Context, product model.Product) (*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
	Update(ctx context.Context, id string, product model.Product) (*model.Product, error)
	Delete(ctx context.Context, id string) error
}

type ProductRepository struct {
	Collection *mongo.Collection
}

// NewProductRepository creates a new repository instance for product-related database operations
func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		Collection: db.Collection("products"),
	}
}

func (r *ProductRepository) Getall(ctx context.Context) *[]model.Product {
	var products []model.Product
	var filter = bson.M{}
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
	}
	err = cur.All(ctx, &products)
	log.Println(products)
	if err != nil {
		log.Println(err)
	}
	return &products
}

// Create inserts a new product into the database
func (r *ProductRepository) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	result, err := r.Collection.InsertOne(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("could not insert product: %v", err)
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.ID = model.MOID(oid.Hex())
	}
	return &product, nil
}

// FindByID retrieves a product by its ID
func (r *ProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("could not find product: %v", err)
	}
	return &product, nil
}

// Update modifies a product's details by ID
func (r *ProductRepository) Update(ctx context.Context, id string, product model.Product) (*model.Product, error) {
	tx, err := r.Collection.Database().Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer tx.EndSession(context.TODO())

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		oid, _ := primitive.ObjectIDFromHex(id)
		update := bson.M{
			"$set": bson.M{
				"name":  product.Name,
				"price": product.Price,
				"stock": product.Stock,
			},
		}
		res, errUpdate := r.Collection.UpdateOne(sessCtx, bson.M{"_id": oid}, update)
		if errUpdate != nil || res.ModifiedCount == 0 {
			tx.AbortTransaction(sessCtx)
			return nil, fmt.Errorf("could not update product: %v", errUpdate)
		}

		if res.MatchedCount > 0 {
			product.ID = model.MOID(oid.Hex())
		}

		if res.ModifiedCount == 1 {
			errUpdate = tx.CommitTransaction(sessCtx)
			if errUpdate != nil {
				log.Println(errUpdate)
				return nil, errUpdate
			}
		}

		return res, nil
	}

	_, err = tx.WithTransaction(context.TODO(), callback)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("could not update product: %v", err)
	}

	return &product, nil
}

// Delete removes a product from the database by ID
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}
	return nil
}
