package middlewares

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "http://domain.com", nil)
	req.Header.Set("X-USER-AUTHENTICATED", "true")

	w := httptest.NewRecorder()
	authCtx := AuthCtx(authHandler)
	authCtx.ServeHTTP(w, req)
	res := w.Body
	body, err := ioutil.ReadAll(res)
	if err != nil {
		t.Error(err)
	}

	if string(body) == `{"code":401,"msg":"Unauthorized"}` {
		t.Error(string(body))
	}
}

func TestAuth_Unauthorized(t *testing.T) {
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "http://domain.com", nil)

	w := httptest.NewRecorder()
	authCtx := AuthCtx(authHandler)
	authCtx.ServeHTTP(w, req)
	res := w.Body
	body, err := ioutil.ReadAll(res)
	if err != nil {
		t.Error(err)
	}

	if string(body) != `{"code":401,"msg":"Unauthorized"}` {
		t.Error(string(body))
	}
}
