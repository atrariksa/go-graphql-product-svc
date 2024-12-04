package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-graphql-product-svc/internal/model"
	"go-graphql-product-svc/internal/service"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/graphql-go/graphql"
)

type ProductHandler struct {
	Service service.IProductService
	cv      service.IClaimsValidator
}

// NewProductHandler creates a new handler for product-related routes
func NewProductHandler(service service.IProductService, cv service.IClaimsValidator) *ProductHandler {
	return &ProductHandler{
		Service: service,
		cv:      cv,
	}
}

// GraphQL Object for Product
var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"name":  &graphql.Field{Type: graphql.String},
		"price": &graphql.Field{Type: graphql.Float},
		"stock": &graphql.Field{Type: graphql.Int},
	},
})

// ServeGraphQL handles GraphQL requests
func (h *ProductHandler) ServeGraphQL(w http.ResponseWriter, r *http.Request) {

	rCtx := r.Context()

	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fields := graphql.Fields{
		"getProduct": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				return h.Service.GetProductByID(p.Context, id)
			},
		},
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return h.Service.GetAllProduct(p.Context), nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	mutationFields := graphql.Fields{
		"createProduct": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"name":  &graphql.ArgumentConfig{Type: graphql.String},
				"price": &graphql.ArgumentConfig{Type: graphql.Float},
				"stock": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name := p.Args["name"].(string)
				price := p.Args["price"].(float64)
				stock := p.Args["stock"].(int)

				product := model.Product{Name: name, Price: price, Stock: stock}
				return h.Service.CreateProduct(p.Context, product)
			},
		},
		"updateProduct": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.String},
				"name":  &graphql.ArgumentConfig{Type: graphql.String},
				"price": &graphql.ArgumentConfig{Type: graphql.Float},
				"stock": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				p.Context = rCtx
				if !h.cv.IsAdmin(rCtx.Value("claims").(jwt.MapClaims)) {
					return nil, errors.New("you cannot access this resource")
				}
				id := p.Args["id"].(string)
				name := p.Args["name"].(string)
				price := p.Args["price"].(float64)
				stock := p.Args["stock"].(int)

				product := model.Product{Name: name, Price: price, Stock: stock}
				return h.Service.UpdateProduct(p.Context, id, product)
			},
		},
		"deleteProduct": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				p.Context = rCtx
				if !h.cv.IsAdmin(rCtx.Value("claims").(jwt.MapClaims)) {
					return nil, errors.New("you cannot access this resource")
				}
				id := p.Args["id"].(string)
				err := h.Service.DeleteProduct(p.Context, id)
				if err != nil {
					return false, err
				}
				return true, nil
			},
		},
	}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: params["query"].(string),
	})

	if len(result.Errors) > 0 {
		http.Error(w, fmt.Sprintf("Failed to execute GraphQL operation: %v", result.Errors), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
