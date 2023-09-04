package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gtngzlv/gophermart/internal/utils"
)

func Authorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := utils.ReturnSecretFromConfig()

		cookie, err := r.Cookie(string(cookieName))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := cookie.Value

		claims := Claims{}
		parsedTokenInfo, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !parsedTokenInfo.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), userID, claims.UserID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
