package inventory

import (
	"errors"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	s := NewStore()
	if _, err := s.Add("SKU-1", "Widget", 5); err != nil {
		t.Fatalf("Add: %v", err)
	}
	it, err := s.Add("SKU-1", "", 3)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if it.Quantity != 8 || it.Name != "Widget" {
		t.Fatalf("got %+v, want qty 8 name Widget", it)
	}
}

func TestAddValidation(t *testing.T) {
	s := NewStore()
	tests := []struct {
		name string
		sku  string
		qty  int
		want error
	}{
		{"empty sku", "", 1, ErrInvalidSKU},
		{"zero qty", "A", 0, ErrInvalidQuantity},
		{"negative qty", "A", -2, ErrInvalidQuantity},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := s.Add(tc.sku, "n", tc.qty); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	s := NewStore()
	_, _ = s.Add("A", "Alpha", 4)

	if _, err := s.Remove("A", 5); !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("got %v, want ErrInsufficientStock", err)
	}
	if _, err := s.Remove("missing", 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
	it, err := s.Remove("A", 4)
	if err != nil || it.Quantity != 0 {
		t.Fatalf("got %+v, %v; want qty 0", it, err)
	}
}

func TestListSorted(t *testing.T) {
	s := NewStore()
	_, _ = s.Add("B", "b", 1)
	_, _ = s.Add("A", "a", 1)
	got := s.List()
	if len(got) != 2 || got[0].SKU != "A" || got[1].SKU != "B" {
		t.Fatalf("unexpected order: %+v", got)
	}
}
