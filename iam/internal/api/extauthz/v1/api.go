package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/dexguitar/spacecraftory/iam/internal/service"
)

const (
	// Cookie/Header names для извлечения session UUID
	SessionCookieName    = "session_uuid"
	SessionCookieNameAlt = "X-Session-Uuid"

	// Headers для передачи session UUID
	HeaderSessionUUID    = "session-uuid"
	HeaderSessionUUIDAlt = "x-session-uuid"
	HeaderAuthorization  = "authorization"
	HeaderCookie         = "cookie"

	// Headers для передачи информации о пользователе в upstream
	HeaderUserUUID    = "x-user-uuid"
	HeaderUserLogin   = "x-user-login"
	HeaderUserEmail   = "x-user-email"
	HeaderSessionExp  = "x-session-expires"
	HeaderContentType = "content-type"
	HeaderAuthStatus  = "x-auth-status"

	ContentTypeJSON  = "application/json"
	AuthStatusDenied = "denied"
)

// API implements envoy.service.auth.v3.Authorization gRPC service
type API struct {
	authv3.UnimplementedAuthorizationServer
	authService service.AuthService
}

// NewAPI creates a new External Authorization API
func NewAPI(authService service.AuthService) *API {
	return &API{
		authService: authService,
	}
}

// Check implements the Envoy ext_authz Check method
func (a *API) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	log.Printf("🔐 External Authorization Check called")

	// Извлекаем session UUID из запроса
	sessionUUID, err := a.extractSessionUUID(req)
	if err != nil {
		log.Printf("❌ Session extraction failed: %v", err)
		return a.denyRequest("Missing or invalid session", typev3.StatusCode_Unauthorized), nil
	}

	log.Printf("📋 Extracted session_uuid: %s", sessionUUID)

	// Проверяем сессию через WhoAmI
	session, user, err := a.authService.WhoAmI(ctx, sessionUUID)
	if err != nil {
		log.Printf("❌ WhoAmI failed: %v", err)
		return a.denyRequest("Invalid session", typev3.StatusCode_Forbidden), nil
	}

	log.Printf("✅ Session valid for user: %s", user.Info.Login)

	// Возвращаем успешный ответ с информацией о пользователе
	return a.allowRequest(user.UUID, user.Info.Login, user.Info.Email, session.ExpiresAt), nil
}

// extractSessionUUID извлекает session UUID из различных источников в запросе
func (a *API) extractSessionUUID(req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil || req.Attributes.Request.Http == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	// 1. Проверяем заголовок session-uuid
	if sessionUUID, ok := headers[HeaderSessionUUID]; ok && sessionUUID != "" {
		return sessionUUID, nil
	}

	// 2. Проверяем заголовок x-session-uuid
	if sessionUUID, ok := headers[HeaderSessionUUIDAlt]; ok && sessionUUID != "" {
		return sessionUUID, nil
	}

	// 3. Проверяем Authorization Bearer token
	if authHeader, ok := headers[HeaderAuthorization]; ok && authHeader != "" {
		sessionUUID := a.extractBearerToken(authHeader)
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	// 4. Проверяем Cookie
	if cookieHeader, ok := headers[HeaderCookie]; ok && cookieHeader != "" {
		sessionUUID := a.extractSessionFromCookies(cookieHeader)
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	return "", fmt.Errorf("session uuid not found in request")
}

// extractBearerToken извлекает токен из заголовка Authorization
func (a *API) extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}

// extractSessionFromCookies извлекает session UUID из cookies
func (a *API) extractSessionFromCookies(cookieHeader string) string {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Add("Cookie", cookieHeader)

	// Пробуем session_uuid
	if cookie, err := req.Cookie(SessionCookieName); err == nil {
		sessionUUID, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			return cookie.Value
		}
		return sessionUUID
	}

	// Пробуем X-Session-Uuid
	if cookie, err := req.Cookie(SessionCookieNameAlt); err == nil {
		sessionUUID, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			return cookie.Value
		}
		return sessionUUID
	}

	return ""
}

// allowRequest создает успешный ответ с заголовками пользователя
func (a *API) allowRequest(userUUID, userLogin, userEmail string, expiresAt time.Time) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserUUID,
				Value: userUUID,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserLogin,
				Value: userLogin,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserEmail,
				Value: userEmail,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderSessionExp,
				Value: expiresAt.Format(time.RFC3339),
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: 0}, // OK
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: headers,
				// Удаляем sensitive заголовки из upstream запроса
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}

// denyRequest создает ответ с отказом в доступе
func (a *API) denyRequest(message string, statusCode typev3.StatusCode) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: statusCode,
				},
				Body: fmt.Sprintf(`{"error": "%s", "timestamp": "%s"}`,
					message, time.Now().Format(time.RFC3339)),
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderContentType,
							Value: ContentTypeJSON,
						},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderAuthStatus,
							Value: AuthStatusDenied,
						},
					},
				},
			},
		},
	}
}
