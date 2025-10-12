package interceptor

import (
	"context"
	"log"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// validator is an interface for requests that have validation
type validator interface {
	Validate() error
}

// ValidationInterceptor creates a server-side unary interceptor that validates
// incoming requests and logs execution time for gRPC methods.
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Extract the method name from the full path
		method := path.Base(info.FullMethod)

		// Log the start of the method
		log.Printf("🚀 Started gRPC method %s\n", method)

		// Validate the request if it implements the validator interface
		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				log.Printf("❌ Validation failed for %s: %v\n", method, err)
				return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
			}
			log.Printf("✅ Validation passed for %s\n", method)
		}

		// Start the timer
		startTime := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Calculate execution duration
		duration := time.Since(startTime)

		// Format message based on result
		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration)
		} else {
			log.Printf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration)
		}

		return resp, err
	}
}
