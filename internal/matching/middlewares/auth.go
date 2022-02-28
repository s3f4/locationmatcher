package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/s3f4/locationmatcher/pkg/apihelper"
)

// AuthCtx gets user from token
func AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := tokenFromHeader(r)

		if token == "" {
			apihelper.Send401(w)
			return
		}

		if err := verifyToken(token); err != nil {
			apihelper.Send401(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenStr string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["authenticated"] == true {
			return nil
		}
	} else {
		return err
	}

	return nil
}

// tokenFromHeader ...
func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
