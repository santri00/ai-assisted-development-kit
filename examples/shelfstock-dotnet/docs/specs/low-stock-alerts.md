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
- Each item has a `threshold` (integer >= 0). Default 0 means "never flagged".
- An item is **low** when `threshold > 0 && quantity <= threshold`.

| Input | Expected |
|---|---|
| `POST /items/A/threshold {"threshold":5}` on existing item | 200, item returned with threshold 5 |
| same on unknown SKU | 404 |
| `{"threshold":-1}` | 400 |
| `GET /items/low-stock` with A(qty 3, thr 5), B(qty 9, thr 5), C(qty 1, thr 0) | 200, `[A]` only, sorted by SKU |

## 4. Interfaces
- `Item` gains `int Threshold` (record gets a new positional parameter; update all constructions).
- `InventoryStore.SetThreshold(string sku, int threshold)`; new `InvalidThresholdException`.
- `InventoryStore.LowStock()` returns `IReadOnlyList<Item>`.
- Routes: `POST /items/{sku}/threshold`, `GET /items/low-stock`.
- Existing routes are unchanged; JSON gains one additive field (`threshold`).

## 5. Edge cases and failure modes
- `GET /items/low-stock` must not be captured by `GET /items/{sku}` (literal segments win; prove it with a test).
- Empty result must serialise as `[]`.
- `Add` and `Remove` must preserve an item's existing threshold (records are immutable: use `with`).
- Threshold changes must be race-free.

## 6. Non-functional
No new packages; O(n) scan is acceptable; no PII.

## 7. Acceptance criteria
- [ ] Given an existing item, When threshold is set to 5, Then `Get` returns threshold 5.
- [ ] Given unknown SKU, When setting threshold, Then `ItemNotFoundException`.
- [ ] Given threshold -1, Then `InvalidThresholdException` (HTTP 400).
- [ ] Given the items in the table above, When `LowStock` runs, Then only A is returned.
- [ ] Adding stock to an item keeps its threshold.
- [ ] `task check` passes.

## 8. Test plan
xUnit tests in `InventoryStoreTests`; API tests in `ApiTests` for 200/400/404 and route precedence.

## 9. Rollout
Additive change; no migration.

## 10. Open questions / assumptions
- Assumption: threshold 0 means disabled (confirmed by owner).
