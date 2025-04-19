package gateway

import (
	"context"
	pb "github.com/kweheliye/omsv2/common/api"
	"github.com/kweheliye/omsv2/common/discovery"
	"log"
)

type stockGRPC struct {
	registry discovery.Registry
}

func NewStockGateway(registry discovery.Registry) *stockGRPC {
	return &stockGRPC{registry}
}

func (g *stockGRPC) CheckIfItemIsInStock(ctx context.Context, customerID string, items []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "stock", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	c := pb.NewStockServiceClient(conn)

	res, err := c.CheckIfItemIsInStock(ctx, &pb.CheckIfItemIsInStockRequest{
		Items: items,
	})

	return res.InStock, res.Items, err
}
