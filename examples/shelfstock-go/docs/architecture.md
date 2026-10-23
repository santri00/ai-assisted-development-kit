# Architecture - shelfstock

## 1. Purpose and context
shelfstock is an in-memory inventory service. Clients (a warehouse UI or scripts) call a JSON API to
add stock, remove stock, and read levels. There is no persistence yet: a restart loses data.

## 2. Building blocks
| Component | Responsibility | Location | Depends on |
|---|---|---|---|
| `main` | Reads config, builds server | `cmd/shelfstock` | httpapi, inventory |
| `httpapi` | Routing, JSON decode/encode, error -> HTTP status | `internal/httpapi` | inventory |
| `inventory` | Domain rules and thread-safe `Store` | `internal/inventory` | stdlib only |

## 3. Key flow: remove stock
1. `POST /items/{sku}/remove {"quantity":N}` reaches `httpapi.remove`.
2. Body is decoded (unknown fields rejected, 1 MiB cap).
3. `Store.Remove` takes the write lock, validates quantity and availability, decrements.
4. Sentinel errors map to 400 / 404 / 409 in `writeErr`; success returns the updated item.

## 4. Data
`Item{SKU, Name, Quantity}` held in a `map[string]Item`. No PII. Stock must never be negative.

## 5. Cross-cutting rules
- Layering: transport -> domain only. `inventory` is free of `net/http`.
- Errors: sentinel errors + `%w` wrapping; status mapping lives only in `httpapi`.
- Concurrency: single `sync.RWMutex` in `Store`.

## 6. Decisions
- [ADR-0001: In-memory store behind a Store type](adr/0001-in-memory-store.md)

## 7. Known debt
- No persistence, no auth, no pagination on `GET /items`.

## 8. Build and run
- Commands: `Taskfile.yaml` (`task --list`). Definition of done: `task check`.
- Container: multi-stage `Dockerfile` (`test` stage = vet + race tests, `runtime` stage = non-root alpine image). Local stack: `compose.yaml` (`task docker:up`).
- Windows: line endings are forced to LF via `.gitattributes`; the race detector runs in Docker (`task docker:check`) because it needs a C compiler.
