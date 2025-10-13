package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"
)

// GetCustomerRequest represents a request to get customer details
type GetCustomerRequest struct {
	CustomerID string `json:"customer_id"`
}

// GetCustomerResponse represents a response with customer details
type GetCustomerResponse struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

// ValidateCustomerRequest represents a request to validate customer
type ValidateCustomerRequest struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
}

// ValidateCustomerResponse represents a validation response
type ValidateCustomerResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// RPCServerExample demonstrates RPC server usage for customer operations
func RPCServerExample(client rabbitmq.Client) error {
	ctx := context.Background()

	// Create RPC server for getting customer details
	server, err := rpc.NewServer(client, rpc.ServerConfig{
		QueueName: "customer.get.rpc",
	})
	if err != nil {
		return fmt.Errorf("failed to create RPC server: %w", err)
	}
	defer server.Close()

	// Start serving RPC requests
	go func() {
		err := server.ServeJSON(ctx, func(ctx context.Context, request interface{}, headers map[string]interface{}) (interface{}, error) {
			req, ok := request.(*GetCustomerRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type")
			}

			log.Printf("RPC Server: Received request for customer %s", req.CustomerID)

			// Simulate database lookup
			time.Sleep(100 * time.Millisecond)

			// Return customer details
			response := GetCustomerResponse{
				CustomerID: req.CustomerID,
				Email:      "customer@example.com",
				Name:       "Customer Name",
				Status:     "active",
			}

			return response, nil
		}, &GetCustomerRequest{})

		if err != nil {
			log.Printf("Error serving RPC requests: %v", err)
		}
	}()

	log.Println("RPC Server: Started serving customer.get.rpc")
	return nil
}

// RPCClientExample demonstrates RPC client usage for customer operations
func RPCClientExample(client rabbitmq.Client) error {
	ctx := context.Background()

	// Create RPC client
	rpcClient, err := rpc.NewClient(client, rpc.ClientConfig{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}
	defer rpcClient.Close()

	// Make RPC call to get customer details
	request := GetCustomerRequest{
		CustomerID: "cust-789",
	}

	var response GetCustomerResponse
	if err := rpcClient.CallJSON(ctx, "customer.get.rpc", request, &response); err != nil {
		return fmt.Errorf("RPC call failed: %w", err)
	}

	log.Printf("RPC Client: Received response: %+v", response)
	return nil
}

// RPCValidationServerExample demonstrates RPC server for customer validation
func RPCValidationServerExample(client rabbitmq.Client) error {
	ctx := context.Background()

	// Create RPC server for customer validation
	server, err := rpc.NewServer(client, rpc.ServerConfig{
		QueueName: "customer.validate.rpc",
	})
	if err != nil {
		return fmt.Errorf("failed to create RPC server: %w", err)
	}
	defer server.Close()

	// Start serving validation requests
	go func() {
		err := server.ServeJSON(ctx, func(ctx context.Context, request interface{}, headers map[string]interface{}) (interface{}, error) {
			req, ok := request.(*ValidateCustomerRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type")
			}

			log.Printf("RPC Server: Validating customer %s with email %s", req.CustomerID, req.Email)

			// Simulate validation logic
			time.Sleep(50 * time.Millisecond)

			// Return validation result
			response := ValidateCustomerResponse{
				Valid:   true,
				Message: "Customer is valid",
			}

			return response, nil
		}, &ValidateCustomerRequest{})

		if err != nil {
			log.Printf("Error serving validation requests: %v", err)
		}
	}()

	log.Println("RPC Server: Started serving customer.validate.rpc")
	return nil
}
