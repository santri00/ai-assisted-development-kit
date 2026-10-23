# ADR-0001: In-memory store behind an InventoryStore type

- Status: Accepted
- Date: 2026-10-21
- Deciders: Sanat (engineering)

## Context
We need a working inventory service quickly to practise AI-assisted workflows. Real persistence would add
packages and operational concerns that distract from the exercise.

## Decision
We will keep inventory in a thread-safe in-memory `InventoryStore` registered as a singleton. All access goes
through its methods so persistence can be introduced later behind an interface without changing endpoints.

## Options considered
1. **In-memory `InventoryStore`** - simple, no packages; data lost on restart.
2. **EF Core + SQLite** - durable; adds packages, migrations, and a DbContext lifetime to reason about.
3. **Do nothing / mock only** - cannot demonstrate real behaviour.

## Consequences
- Positive: fast tests, no external services, no NuGet dependencies beyond test tooling.
- Negative: no durability; single-process only (also: a singleton store means API tests must create a fresh factory per test).
- Follow-up: introduce an `IInventoryRepository` when persistence is required (new ADR).

## Notes for AI agents
- Never expose or mutate the internal dictionary outside `InventoryStore`.
- Do not add EF Core or any database package without a new ADR.
