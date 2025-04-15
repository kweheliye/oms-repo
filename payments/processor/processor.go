package processor

import pb "github.com/kweheliye/omsv2/common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
