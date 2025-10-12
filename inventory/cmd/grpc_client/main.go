package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

const serverAddress = "localhost:50051"

// listParts lists all parts
func listParts(ctx context.Context, client inventoryV1.InventoryServiceClient, filters *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error) {
	resp, err := client.ListParts(ctx, &inventoryV1.ListPartsRequest{Filter: filters})
	if err != nil {
		return nil, err
	}

	return resp.Parts, nil
}

// getPart gets a part by UUID
func getPart(ctx context.Context, client inventoryV1.InventoryServiceClient, uuid string) (*inventoryV1.Part, error) {
	resp, err := client.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: uuid})
	if err != nil {
		return nil, err
	}

	return resp.Part, nil
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
	client := inventoryV1.NewInventoryServiceClient(conn)

	log.Println("=== Testing Inventory API ===")
	log.Println()

	// 1. List all parts with no filters
	log.Println("🪛 Fetching parts with no filters")
	log.Println("===========================")
	parts, err := listParts(ctx, client, nil)
	if err != nil {
		log.Printf("Error listing parts: %v\n", err)
		return
	}

	// Print parts
	log.Printf("Fetched %d parts: %v\n", len(parts), parts)

	// 2. Get each part's information by uuid
	log.Println("🔍 Getting each part's information by uuid")
	log.Println("==================================")

	for i, part := range parts {
		part, err := getPart(ctx, client, part.Uuid)
		if err != nil {
			log.Printf("Error getting part: %v\n", err)
			return
		}

		log.Printf("Fetched part %d: %v\n", i, part.Name)
	}

	log.Println("Testing completed!")
}
