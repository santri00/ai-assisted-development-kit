// Package httpapi exposes the inventory Store over a small JSON HTTP API.
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/shelfstock/internal/inventory"
)

// New returns an http.Handler serving the inventory API.
//
//	GET  /healthz          liveness probe
//	GET  /items            list items
//	GET  /items/{sku}      get one item
//	POST /items/{sku}/add  body {"name":"...","quantity":N}
//	POST /items/{sku}/remove body {"quantity":N}
func New(store *inventory.Store) http.Handler {
	h := &handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /items", h.list)
	mux.HandleFunc("GET /items/{sku}", h.get)
	mux.HandleFunc("POST /items/{sku}/add", h.add)
	mux.HandleFunc("POST /items/{sku}/remove", h.remove)
	return mux
}

type handler struct{ store *inventory.Store }

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type qtyBody struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func (h *handler) list(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	it, err := h.store.Get(r.PathValue("sku"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (h *handler) add(w http.ResponseWriter, r *http.Request) {
	var b qtyBody
	if !decode(w, r, &b) {
		return
	}
	it, err := h.store.Add(r.PathValue("sku"), b.Name, b.Quantity)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (h *handler) remove(w http.ResponseWriter, r *http.Request) {
	var b qtyBody
	if !decode(w, r, &b) {
		return
	}
	it, err := h.store.Remove(r.PathValue("sku"), b.Quantity)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return false
	}
	return true
}

func writeErr(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, inventory.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, inventory.ErrInsufficientStock):
		status = http.StatusConflict
	case errors.Is(err, inventory.ErrInvalidQuantity), errors.Is(err, inventory.ErrInvalidSKU):
		status = http.StatusBadRequest
	}
	writeJSON(w, status, map[string]string{"error": strings.TrimSpace(err.Error())})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
