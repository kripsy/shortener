package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	//nolint:depguard
	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// this middleware for try get jwt from cookie.
//  1. if URL not protected - generate/update new jwt and set in cookie or pass if jwt is valid
//  2. if URL is protected:
//     2.1. if jwt valid and URL protected - pass
//     2.2. if jwt invalid/empty - generate/update new jwt

// JWTMiddleware implement auth functional in hendlers.
// If URL not protected - generate/update new jwt and set in cookie or pass if jwt is valid.
// if URL protected:
//  1. if jwt valid and URL protected - pass
//  2. if jwt invalid/empty - generate/update new jwt
func (m *MyMiddleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protectedURL := []string{
			"/api/user/urls",
		}
		m.MyLogger.Debug("Start JWTMiddleware")

		// check if current URL is protected
		isURLProtected := utils.StingContains(protectedURL, r.URL.Path)
		m.MyLogger.Debug("URL protected value", zap.Bool("msg", isURLProtected))

		// try get token from header

		token, err := utils.GetToken(r)

		// if token empty and url is protected -  return 401
		if err != nil {
			fmt.Printf("Error split bearer token %s", err.Error())
			m.MyLogger.Debug("Error split bearer token", zap.String("msg", err.Error()))
			m.MyLogger.Debug("Create new token")
			//nolint:contextcheck
			_, err = m.setNewCookie(context.Background(), w, r)
			if err != nil {
				m.MyLogger.Debug("Error set cookie", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
			next.ServeHTTP(w, r)

			return
		}

		tokenIsValid, _ := auth.IsTokenValid(token)
		if !tokenIsValid {
			//nolint:contextcheck
			token, err = m.setNewCookie(context.Background(), w, r)
			if err != nil {
				m.MyLogger.Debug("Error set cookie", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}
		m.MyLogger.Debug("tokenIsValid", zap.Bool("msg", tokenIsValid))

		_, err = auth.GetUserID(token)
		if err != nil {
			if isURLProtected {
				m.MyLogger.Debug("Error get user", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)

				return
			}
			// if url not protected - create new token
			m.MyLogger.Debug("Create new token")
			//nolint:contextcheck
			_, err = m.setNewCookie(context.Background(), w, r)
			if err != nil {
				m.MyLogger.Debug("Error set cookie", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
			next.ServeHTTP(w, r)

			return
		}
		next.ServeHTTP(w, r)
	})
}

// setNewCookie create new cookie into cookie, request and writer Headers.
func (m *MyMiddleware) setNewCookie(_ context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	// generate new token
	//nolint:contextcheck
	newUser, err := m.repo.RegisterUser(context.Background())
	if err != nil {
		m.MyLogger.Debug("Error RegisterUser in setNewCookie", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}
	m.MyLogger.Debug("Created new User", zap.Any("User:", newUser))
	token, err := auth.BuildJWTString(newUser.ID)
	if err != nil {
		m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	m.MyLogger.Debug("Token was generated", zap.String("msg", token))
	// get expired time token for set in cookie
	expTime, err := auth.GetExpires(token)
	if err != nil {
		m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	// generate cookie
	//nolint:exhaustruct
	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	}
	w.Header().Add("Authorization", "Bearer "+token)
	r.Header.Add("Authorization", "Bearer "+token)
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)

	return token, nil
}

func (m *MyMiddleware) GrpcJWTInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	m.MyLogger.Debug("Start JWTInterceptor")

	protectedMethods := []string{
		"/Shortener/MethodName", // Замените на имя вашего сервиса и метода
	}

	isMethodProtected := utils.StingContains(protectedMethods, info.FullMethod)
	m.MyLogger.Debug("Method protected value", zap.Bool("msg", isMethodProtected))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "metadata is not provided"))
	}

	values := md["authorization"]
	if len(values) == 0 {
		ctx, err := m.handleUnauthenticated(ctx, isMethodProtected)
		if err != nil {
			return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "internal error to create token"))
		}

		return handler(ctx, req)
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	tokenIsValid, _ := auth.IsTokenValid(token)
	if !tokenIsValid {
		return m.handleUnauthenticated(ctx, isMethodProtected)
	}

	m.MyLogger.Debug("tokenIsValid", zap.Bool("msg", tokenIsValid))

	_, err := auth.GetUserID(token)
	if err != nil && isMethodProtected {
		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "invalid user"))
	}

	return handler(ctx, req)
}

func (m *MyMiddleware) handleUnauthenticated(ctx context.Context, isMethodProtected bool) (context.Context, error) {
	if isMethodProtected {
		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "authentication required"))
	}

	// Generate new token and set it in metadata for further processing
	newToken, err := m.generateNewToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "failed to generate token"))
	}

	newMD := metadata.Pairs("authorization", "Bearer "+newToken)
	newCtx := metadata.NewOutgoingContext(ctx, newMD)
	err = grpc.SendHeader(ctx, newMD)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "failed to set token"))
	}

	return newCtx, nil
}

func (m *MyMiddleware) generateNewToken(ctx context.Context) (string, error) {
	newUser, err := m.repo.RegisterUser(ctx)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	//nolint:wrapcheck
	return auth.BuildJWTString(newUser.ID)
}
