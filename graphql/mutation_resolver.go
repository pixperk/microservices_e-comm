package main

import (
	"context"
	"errors"
	"time"

	"github.com/pixperk/microservices_e-comm/order"
)

var (
	ErrInvalidProductInput = errors.New("invalid product input")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(ctx, input.Name)

	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   a.ID,
		Name: a.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {

	p, err := r.server.catalogClient.PostProduct(input.Name, input.Description, input.Price)

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	products := []order.OrderedProduct{}
	for _, p := range input.Products {
		if p.ID == "" || p.Quantity <= 0 {
			return nil, ErrInvalidProductInput
		}

		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: p.Quantity})
	}
	o, err := r.server.orderClient.PostOrder(ctx, input.AccountID, products)
	if err != nil {
		return nil, err
	}
	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
	}, nil
}
