# Spec: Low-stock threshold and report

- Status: Approved (worked example for the handbook)
- Owner: Sanat
- Links: docs/architecture.md, ADR-0001

## 1. Problem and goal
Warehouse staff need to see which items are running low so they can reorder. Success: a client can set a
reorder threshold per SKU and fetch all items at or below their threshold.

## 2. Scope
**In scope:** per-SKU threshold, `GET /items/low-stock`, validation, tests.
**Out of scope:** notifications (email/Slack), default global threshold, persistence, authentication.

## 3. Behaviour
- Each item has an optional `threshold` (integer >= 0). Default 0 means "never flagged".
- An item is **low** when `threshold > 0 && quantity <= threshold`.

| Input | Expected |
|---|---|
| `POST /items/A/threshold {"threshold":5}` on existing item | 200, item returned with threshold 5 |
| same on unknown SKU | 404 |
| `{"threshold":-1}` | 400 |
| `GET /items/low-stock` with A(qty 3, thr 5), B(qty 9, thr 5), C(qty 0, thr 0) | 200, `[A]` only, sorted by SKU |

## 4. Interfaces
- `Item` gains `Threshold int \`json:"threshold"\``.
- `Store.SetThreshold(sku string, threshold int) (Item, error)`; new sentinel `ErrInvalidThreshold`.
- `Store.LowStock() []Item`.
- Routes: `POST /items/{sku}/threshold`, `GET /items/low-stock`.
- Existing routes and JSON fields are unchanged (additive only).

## 5. Edge cases and failure modes
- Route `GET /items/low-stock` must not be captured by `GET /items/{sku}` (Go 1.22 mux prefers the more specific literal).
- Threshold changes must be race-free.

## 6. Non-functional
Standard library only; O(n) scan is acceptable; no PII.

## 7. Acceptance criteria
- [ ] Given an existing item, When threshold is set to 5, Then `Get` returns threshold 5.
- [ ] Given unknown SKU, When setting threshold, Then `ErrNotFound`.
- [ ] Given threshold -1, Then `ErrInvalidThreshold` (HTTP 400).
- [ ] Given items from the table above, When `LowStock` runs, Then only A is returned.
- [ ] `task check` passes, including `-race`.

## 8. Test plan
Table-driven unit tests in `inventory`; handler tests in `httpapi` for 200/400/404 and route precedence.

## 9. Rollout
Additive change; no migration.

## 10. Open questions / assumptions
- Assumption: threshold 0 means disabled (confirmed by owner).
