package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/shelfstock/internal/inventory"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestAddThenGet(t *testing.T) {
	h := New(inventory.NewStore())
	if rec := do(t, h, "POST", "/items/A/add", `{"name":"Alpha","quantity":3}`); rec.Code != http.StatusOK {
		t.Fatalf("add: %d %s", rec.Code, rec.Body)
	}
	rec := do(t, h, "GET", "/items/A", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"quantity":3`) {
		t.Fatalf("get: %d %s", rec.Code, rec.Body)
	}
}

func TestStatusCodes(t *testing.T) {
	h := New(inventory.NewStore())
	do(t, h, "POST", "/items/A/add", `{"quantity":1}`)
	tests := []struct {
		name         string
		method, path string
		body         string
		want         int
	}{
		{"missing item", "GET", "/items/nope", "", http.StatusNotFound},
		{"bad json", "POST", "/items/A/add", `{`, http.StatusBadRequest},
		{"zero qty", "POST", "/items/A/add", `{"quantity":0}`, http.StatusBadRequest},
		{"oversell", "POST", "/items/A/remove", `{"quantity":9}`, http.StatusConflict},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if rec := do(t, h, tc.method, tc.path, tc.body); rec.Code != tc.want {
				t.Fatalf("got %d, want %d (%s)", rec.Code, tc.want, rec.Body)
			}
		})
	}
}

func TestHealthz(t *testing.T) {
	rec := do(t, New(inventory.NewStore()), "GET", "/healthz", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("healthz: %d %s", rec.Code, rec.Body)
	}
}
