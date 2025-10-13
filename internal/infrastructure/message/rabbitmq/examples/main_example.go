package examples

import (
	"context"
	"log"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"
)

// RunAllExamples demonstrates all RabbitMQ features
func RunAllExamples() {
	// Initialize RabbitMQ client from Viper config
	client, err := rabbitmq.NewClientFromViper()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	log.Println("=== RabbitMQ Examples Started ===")

	// Run Direct Exchange Example
	log.Println("\n--- Direct Exchange Example ---")
	if err := DirectExchangeExample(client); err != nil {
		log.Printf("Direct exchange example failed: %v", err)
	}

	// Wait a bit for messages to be processed
	time.Sleep(2 * time.Second)

	// Run Topic Exchange Example
	log.Println("\n--- Topic Exchange Example ---")
	if err := TopicExchangeExample(client); err != nil {
		log.Printf("Topic exchange example failed: %v", err)
	}

	// Wait a bit for messages to be processed
	time.Sleep(2 * time.Second)

	// Run RPC Server Example
	log.Println("\n--- RPC Server Example ---")
	if err := RPCServerExample(client); err != nil {
		log.Printf("RPC server example failed: %v", err)
	}

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Run RPC Client Example
	log.Println("\n--- RPC Client Example ---")
	if err := RPCClientExample(client); err != nil {
		log.Printf("RPC client example failed: %v", err)
	}

	// Run RPC Validation Server Example
	log.Println("\n--- RPC Validation Server Example ---")
	if err := RPCValidationServerExample(client); err != nil {
		log.Printf("RPC validation server example failed: %v", err)
	}

	log.Println("\n=== All Examples Completed ===")
	log.Println("Press Ctrl+C to exit...")

	// Keep the program running to process messages
	select {}
}

// RunDirectExchangeOnly runs only the direct exchange example
func RunDirectExchangeOnly() {
	client, err := rabbitmq.NewClientFromViper()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	log.Println("=== Direct Exchange Example ===")
	if err := DirectExchangeExample(client); err != nil {
		log.Fatalf("Direct exchange example failed: %v", err)
	}

	log.Println("Press Ctrl+C to exit...")
	select {}
}

// RunTopicExchangeOnly runs only the topic exchange example
func RunTopicExchangeOnly() {
	client, err := rabbitmq.NewClientFromViper()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	log.Println("=== Topic Exchange Example ===")
	if err := TopicExchangeExample(client); err != nil {
		log.Fatalf("Topic exchange example failed: %v", err)
	}

	log.Println("Press Ctrl+C to exit...")
	select {}
}

// RunRPCOnly runs only the RPC example
func RunRPCOnly() {
	client, err := rabbitmq.NewClientFromViper()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	log.Println("=== RPC Example ===")

	// Start RPC servers
	if err := RPCServerExample(client); err != nil {
		log.Fatalf("RPC server example failed: %v", err)
	}

	if err := RPCValidationServerExample(client); err != nil {
		log.Fatalf("RPC validation server example failed: %v", err)
	}

	// Wait for servers to start
	time.Sleep(1 * time.Second)

	// Make RPC calls
	if err := RPCClientExample(client); err != nil {
		log.Printf("RPC client example failed: %v", err)
	}

	// Make validation call
	log.Println("\n--- Making validation RPC call ---")
	rpcClient, err := rpc.NewClient(client, rpc.ClientConfig{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create RPC client: %v", err)
	}

	request := ValidateCustomerRequest{
		CustomerID: "cust-999",
		Email:      "test@example.com",
	}

	var response ValidateCustomerResponse
	if err := rpcClient.CallJSON(ctx, "customer.validate.rpc", request, &response); err != nil {
		log.Printf("Validation RPC call failed: %v", err)
	} else {
		log.Printf("Validation response: %+v", response)
	}

	log.Println("\nPress Ctrl+C to exit...")
	select {}
}
