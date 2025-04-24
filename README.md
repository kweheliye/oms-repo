# OMS (Order Management System)

A microservices-based order management system for food ordering and delivery, built with Go.

## Overview

OMS is a distributed system that manages the entire lifecycle of food orders, from creation to payment processing, kitchen preparation, and delivery. The system is built using a microservices architecture, with each service responsible for a specific domain of the business.

## Architecture

The system consists of the following microservices:

```
┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐
│ Gateway │──────│ Orders  │──────│ Payment │──────│ Kitchen │──────│  Stock  │
│ Service │      │ Service │      │ Service │      │ Service │      │ Service │
└─────────┘      └─────────┘      └─────────┘      └─────────┘      └─────────┘
     │                │                │                │                │
     │                │                │                │                │
     └────────────────┴────────────────┴────────────────┴────────────────┘
                                      │
                                      │
                               ┌──────┴───────┐
                               │   RabbitMQ   │
                               └──────────────┘
```

### Communication Patterns

- **Synchronous Communication**: Services communicate with each other using gRPC for direct, synchronous requests.
- **Asynchronous Communication**: Services publish and subscribe to events using RabbitMQ for asynchronous communication.

### Service Discovery

- **Consul**: Used for service registration, discovery, and health checking.

### Observability

- **Jaeger**: Used for distributed tracing to monitor and troubleshoot transactions across services.

## Components

### Gateway Service

The entry point for client applications. It exposes HTTP endpoints for creating and retrieving orders, and communicates with the Orders service using gRPC.

**Endpoints:**
- `POST /api/customers/{customerID}/orders` - Create a new order
- `GET /api/customers/{customerID}/orders/{orderID}` - Get order details

### Orders Service

Manages order creation, validation, and status updates. It communicates with the Stock service to check item availability and with the Payment service to initiate payment processing.

**Responsibilities:**
- Create and store orders
- Validate orders with the Stock service
- Update order status
- Publish order events to RabbitMQ

### Payments Service

Handles payment processing using Stripe. It creates payment links for orders and processes webhook callbacks from Stripe.

**Responsibilities:**
- Create payment links for orders
- Process payment webhooks
- Update order status after payment
- Publish payment events to RabbitMQ

### Kitchen Service

Processes paid orders for food preparation. It listens for payment events from RabbitMQ and updates order status when food is ready.

**Responsibilities:**
- Listen for paid order events
- Simulate food preparation
- Update order status to "ready"

### Stock Service

Manages inventory and checks item availability. It provides information about items including prices and quantities.

**Responsibilities:**
- Check if items are in stock
- Provide item details including prices
- Update inventory levels

## Technologies Used

- **Go**: Programming language for all services
- **gRPC**: For synchronous inter-service communication
- **RabbitMQ**: For asynchronous messaging between services
- **MongoDB**: For order data storage
- **Consul**: For service discovery and registration
- **Jaeger**: For distributed tracing
- **Stripe**: For payment processing
- **Docker**: For containerization and deployment

## Setup and Installation

### Prerequisites

- Go 1.20+
- Docker and Docker Compose
- Stripe CLI (for payment testing)

### Running with Docker Compose

For external services like MongoDB, RabbitMQ, Consul, and Jaeger, you can use Docker Compose:

```bash
docker compose up --build
```

### Starting the Services

Each service can be started individually using the Go development server or with Air for hot reloading:

```bash
# In separate terminals
cd gateway && air
cd orders && air
cd payments && air
cd kitchen && air
cd stock && air
```

### Setting Up Stripe for Payment Testing

1. Login to Stripe CLI:
```bash
stripe login
```

2. Listen for webhooks:
```bash
stripe listen --forward-to localhost:8081/webhook
```

Where `localhost:8081/webhook` is the endpoint of the payment service's HTTP server.

For testing, use the Stripe test card: `4242424242424242`

## Management UIs

### RabbitMQ Management UI

Access the RabbitMQ management interface at:
```
http://localhost:15672/
```
Default credentials: guest/guest

### Jaeger UI

Access the Jaeger UI for distributed tracing at:
```
http://localhost:16686/
```

### MongoDB Express

Access the MongoDB web interface at:
```
http://localhost:8082/
```


### Pending

Each service should include unit tests for core business logic and integration tests for API endpoints.

