# Architecture - ShelfStock (.NET)

## 1. Purpose and context
ShelfStock is an in-memory inventory service. Clients (a warehouse UI or scripts) call a JSON API to add
stock, remove stock, and read levels. There is no persistence yet: a restart loses data.

## 2. Building blocks
| Component | Responsibility | Location | Depends on |
|---|---|---|---|
| `ShelfStock.Api` | Composition root, routing, JSON, exception -> HTTP status | `src/ShelfStock.Api` | Core |
| `ShelfStock.Core` | Domain rules, thread-safe `InventoryStore`, domain exceptions | `src/ShelfStock.Core` | BCL only |
| `ShelfStock.Tests` | Unit tests (Core) and in-process API tests | `tests/ShelfStock.Tests` | Core, Api |

## 3. Key flow: remove stock
1. `POST /items/{sku}/remove {"quantity":N}` reaches the minimal-API handler in `InventoryEndpoints`.
2. The body binds to `RemoveRequest` (unknown fields rejected, 1 MiB cap, malformed -> 400).
3. `InventoryStore.Remove` takes the lock, validates quantity and availability, decrements.
4. `InventoryException` subclasses map to 400 / 404 / 409 in `StatusFor`; success returns the updated item.

## 4. Data
`Item(Sku, Name, Quantity)` held in a `Dictionary<string, Item>`. No PII. Stock must never be negative.

## 5. Cross-cutting rules
- Layering is enforced by project references: Core cannot see Api or ASP.NET.
- Errors: domain exceptions in Core, status mapping only in Api.
- Concurrency: one `lock` in `InventoryStore`.

## 6. Decisions
- [ADR-0001: In-memory store behind an InventoryStore type](adr/0001-in-memory-store.md)

## 7. Known debt
- No persistence, no auth, no pagination on `GET /items`.

## 8. Build and run
- Commands: `Taskfile.yaml` (`task --list`). Definition of done: `task check`.
- SDK/TFM: `global.json` and `Directory.Build.props`. Analyzer/style strictness: `.editorconfig` + `TreatWarningsAsErrors`.
- Container: multi-stage `Dockerfile` (`test` stage = format check + tests, `runtime` stage = non-root ASP.NET image). Local stack: `compose.yaml` (`task docker:up`).
- Windows: LF line endings are enforced via `.gitattributes` and `.editorconfig`.
