package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	common "github.com/kweheliye/omsv2/common"
	"github.com/kweheliye/omsv2/common/discovery"
	"github.com/kweheliye/omsv2/common/discovery/consul"
	"github.com/kweheliye/omsv2/gateway/gateway"
	"log"
	"net/http"
	"time"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":3000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	jaegerAddr  = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		log.Fatal("failed to set global tracer")
	}

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Minute * 2)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	ordersGateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(ordersGateway)

	mux := http.NewServeMux()
	handler.registerRoutes(mux)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("x to start server: ", err)
	}
}
