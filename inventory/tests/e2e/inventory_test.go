//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Create gRPC client
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "expected successful connection to gRPC application")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Clean up collection after test
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "expected successful cleanup of parts collection")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Insert test part
			var err error
			partUUID, err = env.InsertTestPart(ctx, "Engine V8", "Powerful V8 engine", 5000.0, inventoryV1.Category_CATEGORY_ENGINE)
			Expect(err).ToNot(HaveOccurred(), "expected successful insertion of test part into MongoDB")
		})

		It("should successfully return part by UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(partUUID))
			Expect(resp.GetPart().Name).ToNot(BeEmpty())
			Expect(resp.GetPart().Description).ToNot(BeEmpty())
			Expect(resp.GetPart().Price).To(BeNumerically(">", 0))
			Expect(resp.GetPart().StockQuantity).To(BeNumerically(">=", 0))
			Expect(resp.GetPart().Category).ToNot(Equal(inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED))
		})

		It("should return error for non-existent UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: "non-existent-uuid",
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ListParts", func() {
		BeforeEach(func() {
			// Insert several test parts
			_, err := env.InsertTestPart(ctx, "Engine V8", "Powerful V8 engine", 5000.0, inventoryV1.Category_CATEGORY_ENGINE)
			Expect(err).ToNot(HaveOccurred())

			_, err = env.InsertTestPart(ctx, "Rocket Fuel Tank", "High capacity fuel tank", 3000.0, inventoryV1.Category_CATEGORY_FUEL)
			Expect(err).ToNot(HaveOccurred())

			_, err = env.InsertTestPart(ctx, "Observation Porthole", "Crystal clear porthole", 1500.0, inventoryV1.Category_CATEGORY_PORTHOLE)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return all parts without filter", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: nil,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())
			Expect(len(resp.GetParts())).To(BeNumerically(">=", 3))
		})

		It("should filter parts by category", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())
			for _, part := range resp.GetParts() {
				Expect(part.Category).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))
			}
		})

		It("should filter parts by name", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"Engine V8"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())
			Expect(resp.GetParts()[0].Name).To(Equal("Engine V8"))
		})

		It("should filter parts by UUID", func() {
			// Insert part and get its UUID
			partUUID, err := env.InsertTestPart(ctx, "Special Wing", "Unique wing design", 2500.0, inventoryV1.Category_CATEGORY_WING)
			Expect(err).ToNot(HaveOccurred())

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{partUUID},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(1))
			Expect(resp.GetParts()[0].Uuid).To(Equal(partUUID))
			Expect(resp.GetParts()[0].Name).To(Equal("Special Wing"))
		})

		It("should return empty list for non-existent filters", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"NonExistentPart"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(BeEmpty())
		})
	})
})
