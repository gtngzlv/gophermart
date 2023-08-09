package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gtngzlv/gophermart/internal/utils"
	"net/http"
)

var allowList = map[string]bool{
	"/api/user/register": true,
	"/api/user/login":    true,
}

var (
	cookieName = "newcookie"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func Authorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := utils.ReturnSecretFromConfig()
		if _, ok := allowList[r.URL.Path]; ok {
			h.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(cookieName)
		if err != nil {
			return
		}

		token := cookie.Value

		claim := Claims{}
		parsedTokenInfo, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !parsedTokenInfo.Valid {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				w.WriteHeader(http.StatusUnauthorized)
			}
			w.WriteHeader(http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), cookieName, claim)
		h.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
