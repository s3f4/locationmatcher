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
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjp0cnVlfQ.iG8Rux4vSqgoBvG2OMggjK9Q5QGyvATykfH8qKbJTAs")

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

func TestAuth_Parse_error(t *testing.T) {
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "http://domain.com", nil)
	req.Header.Set("Authorization", "Bearer .eyJhdXRoZW50aWNhdGVkIjp0cnVlfQ.iG8Rux4vSqgoBvG2OMggjK9Q5QGyvATykfH8qKbJTAs")

	w := httptest.NewRecorder()

	authCtx := AuthCtx(authHandler)
	authCtx.ServeHTTP(w, req)
	res := w.Body
	body, err := ioutil.ReadAll(res)
	if err != nil {
		t.Error(err)
	}

	if string(body) != `{"code":401,"msg":"Unauthorized"}` {
		t.Error("it should be unauthorized")
	}
}

func TestAuth_Valid_Token_Error(t *testing.T) {
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "http://domain.com", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkMiI6dHJ1ZX0.zCddOuhMXwE4bHavgzBaGbHzbhIj3XMqrAyQhPNNCWk")

	w := httptest.NewRecorder()

	authCtx := AuthCtx(authHandler)
	authCtx.ServeHTTP(w, req)
	res := w.Body
	body, err := ioutil.ReadAll(res)
	if err != nil {
		t.Error(err)
	}

	if string(body) != `{"code":401,"msg":"Unauthorized"}` {
		t.Error("it should be unauthorized")
	}
}

func TestAuth_No_Bearer(t *testing.T) {
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "http://domain.com", nil)
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjp0cnVlfQ.iG8Rux4vSqgoBvG2OMggjK9Q5QGyvATykfH8qKbJTAs")

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
