package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dexguitar/spacecraftory/inventory/internal/interceptor"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

// inventoryService implements gRPC service for working with spacecraft parts
type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()

	// If no filter provided, return all parts
	if filter == nil {
		return &inventoryV1.ListPartsResponse{
			Parts: s.getAllParts(),
		}, nil
	}

	// If specific UUIDs are requested
	if len(filter.GetUuids()) > 0 {
		return &inventoryV1.ListPartsResponse{
			Parts: s.getPartsByUUIDs(filter.GetUuids()),
		}, nil
	}

	// Filter by other criteria
	return &inventoryV1.ListPartsResponse{
		Parts: s.filterParts(filter),
	}, nil
}

func (s *inventoryService) getAllParts() []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, part := range s.parts {
		parts = append(parts, part)
	}
	return parts
}

func (s *inventoryService) getPartsByUUIDs(uuids []string) []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, 0, len(uuids))
	for _, uuid := range uuids {
		if part, ok := s.parts[uuid]; ok {
			parts = append(parts, part)
		}
	}
	return parts
}

func (s *inventoryService) filterParts(filter *inventoryV1.PartsFilter) []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, 0)
	for _, part := range s.parts {
		if s.partMatchesFilter(part, filter) {
			parts = append(parts, part)
		}
	}
	return parts
}

func (s *inventoryService) partMatchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	return s.matchesNameFilter(part, filter.GetNames()) &&
		s.matchesCategoryFilter(part, filter.GetCategories()) &&
		s.matchesCountryFilter(part, filter.GetManufacturerCountries()) &&
		s.matchesTagFilter(part, filter.GetTags())
}

// matchesNameFilter checks if a part matches the name filter
func (s *inventoryService) matchesNameFilter(part *inventoryV1.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return slices.Contains(names, part.GetName())
}

// matchesCategoryFilter checks if a part matches the category filter
func (s *inventoryService) matchesCategoryFilter(part *inventoryV1.Part, categories []inventoryV1.Category) bool {
	if len(categories) == 0 {
		return true
	}
	return slices.Contains(categories, part.GetCategory())
}

// matchesCountryFilter checks if a part matches the manufacturer country filter
func (s *inventoryService) matchesCountryFilter(part *inventoryV1.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	return slices.Contains(countries, part.GetManufacturer().GetCountry())
}

// matchesTagFilter checks if a part matches the tag filter
func (s *inventoryService) matchesTagFilter(part *inventoryV1.Part, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	return slices.ContainsFunc(part.GetTags(), func(tag string) bool {
		return slices.Contains(filterTags, tag)
	})
}

// initializeMockData creates some initial parts for testing
func initializeMockData() map[string]*inventoryV1.Part {
	mockParts := map[string]*inventoryV1.Part{
		uuid.NewString(): {
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion engine for interstellar travel",
			Price:         150000.00,
			StockQuantity: 5,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 3.5,
				Width:  2.0,
				Height: 2.5,
				Weight: 500.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "SpaceTech Industries",
				Country: "USA",
				Website: "https://spacetech.example.com",
			},
			Tags:      []string{"quantum", "propulsion", "interstellar"},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
		uuid.NewString(): {
			Name:          "Fusion Fuel Cell",
			Description:   "Advanced fusion-based fuel cell for long-duration missions",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 1.2,
				Width:  0.8,
				Height: 1.0,
				Weight: 100.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "Energy Solutions Corp",
				Country: "Germany",
				Website: "https://energysolutions.example.com",
			},
			Tags:      []string{"fusion", "fuel", "efficient"},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
		uuid.NewString(): {
			Name:          "Reinforced Porthole",
			Description:   "Triple-layered reinforced viewing porthole for crew observation",
			Price:         25000.00,
			StockQuantity: 20,
			Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 1.0,
				Width:  1.0,
				Height: 0.3,
				Weight: 50.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "ViewTech Manufacturing",
				Country: "Japan",
				Website: "https://viewtech.example.com",
			},
			Tags:      []string{"observation", "reinforced", "safety"},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
		uuid.NewString(): {
			Name:          "Aerodynamic Wing Panel",
			Description:   "Carbon-fiber composite wing panel for atmospheric re-entry",
			Price:         45000.00,
			StockQuantity: 8,
			Category:      inventoryV1.Category_CATEGORY_WING,
			Dimensions: &inventoryV1.Dimensions{
				Length: 5.0,
				Width:  2.5,
				Height: 0.5,
				Weight: 200.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "AeroDynamics Ltd",
				Country: "UK",
				Website: "https://aerodynamics.example.com",
			},
			Tags:      []string{"aerodynamic", "reentry", "composite"},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}

	return mockParts
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.ValidationInterceptor()),
	)

	// Enable reflection for debugging
	reflection.Register(s)

	// Register service with mock data
	service := &inventoryService{
		parts: initializeMockData(),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	go func() {
		log.Printf("🚀 Inventory gRPC server listening on port %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Inventory gRPC: %v\n", err)
			return
		}
	}()

	// Launch HTTP server with gRPC gateway
	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Create new serve mux for HTTP requests with dial
		mux := runtime.NewServeMux()

		// Create new dial options for HTTP requests
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Register inventory service handler from endpoint
		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway for inventory service: %v\n", err)
			return
		}

		// Create new HTTP server mux
		httpMux := http.NewServeMux()

		// Mount gRPC-Gateway at /api/
		httpMux.Handle("/api/", mux)

		// Serve Swagger JSON
		httpMux.HandleFunc("/apidocs.swagger.json", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📄 Serving swagger JSON: %s", r.URL.Path)
			http.ServeFile(w, r, "shared/pkg/swagger/inventory/v1/inventory.swagger.json")
		})

		// Serve Swagger UI HTML
		httpMux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📚 Serving swagger UI: %s", r.URL.Path)
			http.ServeFile(w, r, "shared/pkg/swagger/swagger-ui.html")
		})

		// Redirect root to Swagger UI
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🔄 Request to: %s", r.URL.Path)
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
		})

		// Create new HTTP server
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Start HTTP server
		log.Printf("🚀 Inventory HTTP server with gRPC-Gateway listening on port %d\n", httpPort)
		log.Printf("📚 Swagger UI available at: http://localhost:%d/swagger-ui.html\n", httpPort)
		log.Printf("📄 Swagger JSON available at: http://localhost:%d/apidocs.swagger.json\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve Inventory HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Inventory servers...")

	// Shutdown HTTP server
	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Inventory HTTP server shutdown error: %v", err)
		}
		log.Println("✅ Inventory HTTP server stopped")
	}

	// Shutdown gRPC server
	s.GracefulStop()
	log.Println("✅ Inventory gRPC server stopped")
}
