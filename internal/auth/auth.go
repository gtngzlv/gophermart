package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gtngzlv/gophermart/internal/utils"
)

func GenerateCookie(w http.ResponseWriter, login string) error {
	secret := utils.ReturnSecretFromConfig()
	expirationTime := &jwt.NumericDate{Time: time.Now().Add(time.Hour)}
	claims := Claims{
		Login: login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = string(cookieName)
	cookie.Value = tokenString
	cookie.Expires = expirationTime.Time
	http.SetCookie(w, cookie)
	return nil
}

func GetUserLoginFromToken(w http.ResponseWriter, r *http.Request) string {
	claim, ok := r.Context().Value(cookieName).(Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return ""
	}
	return claim.Login
}
