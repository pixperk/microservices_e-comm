package main

import (
	"context"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{
			{
				ID:   r.ID,
				Name: r.Name,
			},
		}, nil
	}

	skip, take := uint64(0), uint64(10)
	if pagination != nil {
		skip, take = pagination.bounds()
	}
	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for _, a := range accountList {
		accounts = append(accounts, &Account{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return accounts, nil
}
func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {

	if id != nil {
		p, err := r.server.catalogClient.GetProduct(*id)
		if err != nil {
			return nil, err
		}
		return []*Product{
			{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		}, nil
	}

	skip, take := uint64(0), uint64(10)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	products, err := r.server.catalogClient.GetProducts(skip, take, nil, q)
	if err != nil {
		return nil, err
	}

	result := make([]*Product, len(products))
	for i, p := range products {
		result[i] = &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}

	return result, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}
