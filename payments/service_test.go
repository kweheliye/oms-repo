package main

import (
	"context"
	pb "github.com/kweheliye/omsv2/common/api"
	inmemRegistry "github.com/kweheliye/omsv2/common/discovery/inmem"
	"github.com/kweheliye/omsv2/payments/gateway"
	"github.com/kweheliye/omsv2/payments/processor/inmem"
	"testing"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()

	newGateway := gateway.NewGateway(registry)
	svc := NewService(processor, newGateway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &pb.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
