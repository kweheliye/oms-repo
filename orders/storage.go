package main

import (
	"context"
	"errors"
	pb "github.com/kweheliye/omsv2/common/api"
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
		ID:         id,
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
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
