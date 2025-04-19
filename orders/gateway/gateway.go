package gateway

import (
	"context"
	pb "github.com/kweheliye/omsv2/common/api"
)

type StockGateway interface {
	CheckIfItemIsInStock(ctx context.Context, customerID string, items []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
}
