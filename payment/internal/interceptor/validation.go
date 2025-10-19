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

type validator interface {
	Validate() error
}

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		method := path.Base(info.FullMethod)

		log.Printf("🚀 Started gRPC method %s\n", method)

		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				log.Printf("❌ Validation failed for %s: %v\n", method, err)
				return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
			}
			log.Printf("✅ Validation passed for %s\n", method)
		}

		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration)
		} else {
			log.Printf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration)
		}

		return resp, err
	}
}
