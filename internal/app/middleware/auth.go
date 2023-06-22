package middleware

import (
	"context"
	"errors"
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

func (m *MyMiddleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protectedURL := []string{
			"/api/user/urls",
		}
		m.MyLogger.Debug("Start JWTMiddleware")

		// check if current URL is protected
		isURLProtected := utils.StingContains(protectedURL, r.URL.Path)
		m.MyLogger.Debug("URL protected value", zap.Bool("msg", isURLProtected))

		// get jwt token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// return if error to get cookie
			if !errors.Is(err, http.ErrNoCookie) {
				m.MyLogger.Debug("Error get Cookie", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if isURLProtected {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			m.MyLogger.Debug("No Cookie")
			m.setNewCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}

		// continue if cookie not empty
		tokenString := cookie.Value
		m.MyLogger.Debug("Current Cookie", zap.String("msg", tokenString))

		tokenIsValid, _ := auth.IsTokenValid(tokenString)
		m.MyLogger.Debug("Current token is valid?", zap.Bool("msg", tokenIsValid))

		if !tokenIsValid {
			err = m.setNewCookie(w, r)
			if err != nil {
				m.MyLogger.Debug("Error set cookie", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		_, err = auth.GetUserID(tokenString)
		if err != nil {
			m.MyLogger.Debug("Error get user ID from token", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *MyMiddleware) setNewCookie(w http.ResponseWriter, r *http.Request) error {
	// generate new token
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	newUser, err := m.repo.RegisterUser(ctx)
	if err != nil {
		m.MyLogger.Debug("Error RegisterUser in setNewCookie", zap.String("msg", err.Error()))
		return err
	}
	m.MyLogger.Debug("Created new User", zap.Any("User:", newUser))
	token, err := auth.BuildJWTString(newUser.ID)
	if err != nil {
		m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))
		return err
	}

	m.MyLogger.Debug("Token was generated", zap.String("msg", token))
	// get expired time token for set in cookie
	expTime, err := auth.GetExpires(token)
	if err != nil {
		m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))
		return err
	}

	// generate cookie
	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	}
	w.Header().Add("Authorization", "Bearer "+token)
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	return nil
}
