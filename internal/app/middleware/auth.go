package middleware

import (
	"net/http"

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
		isURLProtected := utils.StingContains(protectedURL, r.URL.Path)
		m.MyLogger.Debug("URL protected value", zap.Bool("msg", isURLProtected))
		// generate new token
		token, err := auth.BuildJWTString()
		if err != nil {
			m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))
		}

		m.MyLogger.Debug("Token was generated", zap.String("msg", token))
		// get expired time token for set in cookie
		expTime, err := auth.GetExpires(token)
		if err != nil {
			m.MyLogger.Debug("Error JWTMiddleware", zap.String("msg", err.Error()))
		}

		// generate cookie
		cookie := &http.Cookie{
			Name:     "jwtString",
			Value:    token,
			Secure:   true,
			Expires:  expTime,
			SameSite: http.SameSiteDefaultMode,
		}

		// set cookie
		http.SetCookie(w, cookie)

		next.ServeHTTP(w, r)
	})
}
