package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/dexguitar/spacecraftory/platform/pkg/middleware/grpc"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = authV1.AuthServiceClient

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	iamClient IAMClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// client (X-Session-Uuid) -> auth middleware (add session_uuid in ctx (incoming)) -> order api (outgoing)-> auth interceptor ->inventory
// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем session UUID из заголовка
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.WhoAmI(r.Context(), &authV1.WhoAmIRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddUserToContext(ctx, whoamiRes.User)
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
