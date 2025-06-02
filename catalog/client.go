package catalog

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/pixperk/microservices_e-comm/catalog/pb"
)

type Client struct {
	conn          *grpc.ClientConn
	catalogClient pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:          conn,
		catalogClient: pb.NewCatalogServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(name, description string, price float64) (*Product, error) {
	resp, err := c.catalogClient.PostProduct(
		context.Background(),
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
	}, nil
}

func (c *Client) GetProduct(id string) (*Product, error) {
	resp, err := c.catalogClient.GetProduct(
		context.Background(),
		&pb.GetProductRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
	}, nil
}

func (c *Client) GetProducts(skip, take uint64, ids []string, query string) ([]*Product, error) {
	resp, err := c.catalogClient.GetProducts(
		context.Background(),
		&pb.GetProductsRequest{Skip: skip, Take: take, Ids: ids, Query: query},
	)
	if err != nil {
		return nil, err
	}
	products := make([]*Product, len(resp.Products))
	for i, p := range resp.Products {
		products[i] = &Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}
	return products, nil
}
