package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/kweheliye/omsv2/common/api"
	"github.com/kweheliye/omsv2/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
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

func (h *grpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, p)
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	q, eror := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if eror != nil {
		log.Fatal(eror)
	}

	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", q.Name))
	defer messageSpan.End()

	items, err := h.service.ValidateOrder(amqpContext, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(amqpContext, p, items)

	marshallOrder, err := json.Marshal(o)

	if err != nil {
		return nil, err
	}

	// inject the headers
	headers := broker.InjectAMQPHeaders(amqpContext)

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshallOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})
	return o, nil
}
