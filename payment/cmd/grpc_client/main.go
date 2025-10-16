package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const serverAddress = "localhost:50052"

// payOrder pays an order and generates a transaction uuid
func payOrder(ctx context.Context, client paymentV1.PaymentServiceClient, request *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	resp, err := client.PayOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	ctx := context.Background()

	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connection: %v", cerr)
		}
	}()

	// Create gRPC client
	client := paymentV1.NewPaymentServiceClient(conn)

	log.Println("=== Testing Payment API ===")
	log.Println()

	// 1. Pay an order
	log.Println("💳 Paying an order with card")
	log.Println("===========================")
	resp, err := payOrder(ctx, client, &paymentV1.PayOrderRequest{
		OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
		UserUuid:      "550e8400-e29b-41d4-a716-446655440000",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	})
	if err != nil {
		log.Printf("Error paying an order: %v\n", err)
		return
	}
	log.Printf("✅ Order paid with card! Transaction UUID: %v\n", resp.TransactionUuid)

	// 2. Pay an order with SBP
	log.Println("🔂 Paying an order with SBP")
	log.Println("===========================")
	resp, err = payOrder(ctx, client, &paymentV1.PayOrderRequest{
		OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
		UserUuid:      "550e8400-e29b-41d4-a716-446655440000",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
	})
	if err != nil {
		log.Printf("Error paying an order: %v\n", err)
		return
	}
	log.Printf("✅ Order paid with SBP! Transaction UUID: %v\n", resp.TransactionUuid)

	// 3. Pay an order with credit card
	log.Println("💳 Paying an order with credit card")
	log.Println("===========================")
	resp, err = payOrder(ctx, client, &paymentV1.PayOrderRequest{
		OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
		UserUuid:      "550e8400-e29b-41d4-a716-446655440000",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	})
	if err != nil {
		log.Printf("Error paying an order: %v\n", err)
		return
	}
	log.Printf("✅ Order paid with credit card! Transaction UUID: %v\n", resp.TransactionUuid)

	// 4. Pay an order with investor money
	log.Println("🧐 Paying an order with investor money")
	log.Println("===========================")
	resp, err = payOrder(ctx, client, &paymentV1.PayOrderRequest{
		OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
		UserUuid:      "550e8400-e29b-41d4-a716-446655440000",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	})
	if err != nil {
		log.Printf("Error paying an order: %v\n", err)
		return
	}
	log.Printf("✅ Order paid with investor money! Transaction UUID: %v\n", resp.TransactionUuid)

	log.Println("✅✅ Testing completed! ✅✅")
}
