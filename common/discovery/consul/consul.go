package consul

import (
	"context"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"strconv"
	"strings"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	log.Printf("Registering service %s", serviceName)

	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid host:port format: %s", hostPort)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceID,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "3m",
			Timeout:                        "3m",
			DeregisterCriticalServiceAfter: "3m",
		},
	})
}
func (r *Registry) Deregister(ctx context.Context, instanceID, serverName string) error {
	log.Printf("Deregistering service %s", serverName)
	return r.client.Agent().CheckDeregister(serverName)
}

func (r *Registry) HealthCheck(instanceID, serviceName string) error {
	log.Printf("Health Checking service %s", serviceName)
	return r.client.Agent().UpdateTTL(instanceID, "online", consul.HealthPassing)
}

func (r *Registry) Discover(ctx context.Context, serverName string) ([]string, error) {
	log.Printf("Discovery service %s", serverName)
	entries, _, err := r.client.Health().Service(serverName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, entry := range entries {
		instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}

	return instances, nil
}
