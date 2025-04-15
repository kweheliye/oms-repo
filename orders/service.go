package main

import (
	"context"
	"github.com/kweheliye/omsv2/common"
	pb "github.com/kweheliye/omsv2/common/api"
	"log"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store: store}
}

func (s *service) GetOrder(ctx context.Context, pb *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, pb.OrderID, pb.CustomerID)
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id, err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}
	o := &pb.Order{
		ID:         id,
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
	}
	return o, nil

}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Println("merging items", mergedItems)

	//Tempoarrry :

	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1RDmQE4WTCZyoKdfjcuMQUvz",
			Quantity: i.Quantity,
			ID:       i.ID,
		})
	}
	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemWithQuantity) []*pb.ItemWithQuantity {
	merged := make([]*pb.ItemWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, item)
		}
	}

	return merged

}
