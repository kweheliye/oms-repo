package main

import (
	"context"
	pb "github.com/kweheliye/omsv2/common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
