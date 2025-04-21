package gateway

import (
	"context"

	pb "github.com/kweheliye/omsv2/common/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
