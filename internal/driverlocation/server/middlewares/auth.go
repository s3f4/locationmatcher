package middlewares

import (
	"net/http"

	"github.com/s3f4/locationmatcher/pkg/apihelper"
)

func AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-USER-AUTHENTICATED") != "true" {
			apihelper.Send401(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
