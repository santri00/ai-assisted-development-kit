// Command shelfstock runs the inventory HTTP service.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"example.com/shelfstock/internal/httpapi"
	"example.com/shelfstock/internal/inventory"
)

func main() {
	addr := os.Getenv("SHELFSTOCK_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           httpapi.New(inventory.NewStore()),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("shelfstock listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
