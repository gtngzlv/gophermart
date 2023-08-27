package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gtngzlv/gophermart/internal/utils"
)

type (
	cookie string
	login  string
)

const (
	cookieName  cookie = "authToken"
	loginCookie login  = "login"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func returnNewClaims() Claims {
	return Claims{}
}

func GenerateCookie(w http.ResponseWriter, login string) error {
	secret := utils.ReturnSecretFromConfig()
	expirationTime := &jwt.NumericDate{Time: time.Now().Add(time.Hour)}
	claim := returnNewClaims()
	claim.Login = login
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
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
	login, ok := r.Context().Value(loginCookie).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return ""
	}
	return login
}
