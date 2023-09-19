package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gtngzlv/gophermart/internal/utils"
)

type (
	cookie string
	user   string
)

const (
	cookieName cookie = "authToken"
	userID     user   = "userID"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func GenerateCookie(w http.ResponseWriter, userID int) error {
	secret := utils.ReturnSecretFromConfig()
	expirationTime := &jwt.NumericDate{Time: time.Now().Add(time.Hour)}
	claims := Claims{}
	claims.UserID = userID
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

func GetUserIDFromToken(w http.ResponseWriter, r *http.Request) int {
	userID, ok := r.Context().Value(userID).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return 0
	}
	return userID
}
