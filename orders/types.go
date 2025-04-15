package main

import (
	"context"
	pb "github.com/kweheliye/omsv2/common/api"
)

type OrdersService interface {
	CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, pb []*pb.Item) (*pb.Order, error)
	ValidateOrder(ctx context.Context, pb *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(ctx context.Context, pb *pb.GetOrderRequest) (*pb.Order, error)
}

type OrdersStore interface {
	Create(ctx context.Context, p *pb.CreateOrderRequest, pb []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
}
