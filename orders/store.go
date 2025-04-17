package main

import (
	"context"
	"errors"
	pb "github.com/kweheliye/omsv2/common/api"
	"log"
)

var orders = make([]*pb.Order, 0)

type store struct {
	//add mongo db
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	id := "42"
	orders = append(orders, &pb.Order{
		ID:          id,
		CustomerID:  p.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})
	return id, nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {
	for _, o := range orders {
		if o.ID == id && o.CustomerID == customerID {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, o *pb.Order) error {
	log.Println("Received order update", o)
	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = o.Status
			orders[i].PaymentLink = o.PaymentLink
			return nil
		}
	}
	return nil
}
