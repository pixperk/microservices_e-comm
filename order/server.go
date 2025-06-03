package order

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/pixperk/microservices_e-comm/order/pb"

	"github.com/pixperk/microservices_e-comm/account"
	"github.com/pixperk/microservices_e-comm/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		return nil, errors.New("account not found")
	}

	productIDs := []string{}
	orderedProducts, err := s.catalogClient.GetProducts(0, 0, productIDs, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	products := []OrderedProduct{}
	for _, p := range orderedProducts {
		product := &OrderedProduct{
			ID:          p.ID,
			Quantity:    0,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
		}

		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = int(rp.Quantity)
				break
			}
		}
		if product.Quantity > 0 {
			products = append(products, *product)
		}
	}
	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		return nil, fmt.Errorf("failed to post order: %w", err)
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    uint32(p.Quantity),
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for account: %w", err)
	}

	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}
	products, err := s.catalogClient.GetProducts(0, 0, productIDs, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	orders := []*pb.Order{}
	for _, o := range accountOrders {
		op := &pb.Order{
			AccountId:  o.AccountID,
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Decorate orders with products
		for _, product := range o.Products {
			// Populate product fields
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    uint32(product.Quantity),
			})
		}
		orders = append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
