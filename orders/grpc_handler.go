package main

import (
	"context"
	"encoding/json"
	pb "github.com/kweheliye/omsv2/common/api"
	"github.com/kweheliye/omsv2/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService, ch *amqp.Channel) {
	handler := &grpcHandler{service: service, channel: ch}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order received! %v", p)

	items, err := h.service.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(ctx, p, items)

	marshallOrder, err := json.Marshal(o)

	if err != nil {
		return nil, err
	}

	q, eror := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if eror != nil {
		log.Fatal(eror)
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshallOrder,
		DeliveryMode: amqp.Persistent,
	})
	return o, nil
}
