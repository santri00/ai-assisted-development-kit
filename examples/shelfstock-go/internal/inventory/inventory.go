// Package inventory holds the core domain logic for tracking stock levels.
package inventory

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Sentinel errors returned by Store operations.
var (
	ErrNotFound          = errors.New("inventory: item not found")
	ErrInsufficientStock = errors.New("inventory: insufficient stock")
	ErrInvalidQuantity   = errors.New("inventory: quantity must be positive")
	ErrInvalidSKU        = errors.New("inventory: sku must not be empty")
)

// Item is a stocked product identified by SKU.
type Item struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// Store is an in-memory, concurrency-safe inventory.
type Store struct {
	mu    sync.RWMutex
	items map[string]Item
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{items: make(map[string]Item)}
}

// Add increases stock for sku by qty, creating the item if needed.
func (s *Store) Add(sku, name string, qty int) (Item, error) {
	if sku == "" {
		return Item{}, ErrInvalidSKU
	}
	if qty <= 0 {
		return Item{}, ErrInvalidQuantity
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	it := s.items[sku]
	it.SKU = sku
	if name != "" {
		it.Name = name
	}
	it.Quantity += qty
	s.items[sku] = it
	return it, nil
}

// Remove decreases stock for sku by qty. It never lets stock go negative.
func (s *Store) Remove(sku string, qty int) (Item, error) {
	if qty <= 0 {
		return Item{}, ErrInvalidQuantity
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[sku]
	if !ok {
		return Item{}, fmt.Errorf("%w: %s", ErrNotFound, sku)
	}
	if it.Quantity < qty {
		return Item{}, fmt.Errorf("%w: have %d, want %d", ErrInsufficientStock, it.Quantity, qty)
	}
	it.Quantity -= qty
	s.items[sku] = it
	return it, nil
}

// Get returns the item for sku.
func (s *Store) Get(sku string) (Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[sku]
	if !ok {
		return Item{}, fmt.Errorf("%w: %s", ErrNotFound, sku)
	}
	return it, nil
}

// List returns all items sorted by SKU.
func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SKU < out[j].SKU })
	return out
}
