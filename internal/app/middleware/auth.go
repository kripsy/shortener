package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
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

		token, err := utils.GetToken(w, r)

		// if token empty and url is protected -  return 401
		if err != nil {
			fmt.Printf("Error split bearer token %s", err.Error())
			m.MyLogger.Debug("Error split bearer token", zap.String("msg", err.Error()))
			// if isURLProtected {
			// 	m.MyLogger.Debug("Error split bearer token and URL protected")
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	return
			// }
			// if url not protected - create new token
			m.MyLogger.Debug("Create new token")
			_, err = m.setNewCookie(w, r)
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
			token, err = m.setNewCookie(w, r)
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
			_, err = m.setNewCookie(w, r)
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
func (m *MyMiddleware) setNewCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	// generate new token
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	newUser, err := m.repo.RegisterUser(ctx)
	if err != nil {
		m.MyLogger.Debug("Error RegisterUser in setNewCookie", zap.String("msg", err.Error()))
		return "", err
	}
	m.MyLogger.Debug("Created new User", zap.Any("User:", newUser))
	token, err := auth.BuildJWTString(newUser.ID)
	if err != nil {
		m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))
		return "", err
	}

	m.MyLogger.Debug("Token was generated", zap.String("msg", token))
	// get expired time token for set in cookie
	expTime, err := auth.GetExpires(token)
	if err != nil {
		m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))
		return "", err
	}

	// generate cookie
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
