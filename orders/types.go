package main

import (
	"context"
	pb "github.com/kweheliye/omsv2/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersService interface {
	CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, pb []*pb.Item) (*pb.Order, error)
	ValidateOrder(ctx context.Context, pb *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(ctx context.Context, pb *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(ctx context.Context, pb *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(ctx context.Context, p *pb.CreateOrderRequest, pb []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
	Update(ctx context.Context, id string, o *pb.Order) error
}

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID  string             `bson:"customerID,omitempty"`
	Status      string             `bson:"status,omitempty"`
	PaymentLink string             `bson:"paymentLink,omitempty"`
	Items       []*pb.Item         `bson:"items,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		ID:          o.ID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
