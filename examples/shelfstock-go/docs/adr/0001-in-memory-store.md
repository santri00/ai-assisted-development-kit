# ADR-0001: In-memory store behind a Store type

- Status: Accepted
- Date: 2026-10-21
- Deciders: Sanat (engineering)

## Context
We need a working inventory service quickly to practise AI-assisted workflows. Real persistence would add
dependencies and operational concerns that distract from the exercise.

## Decision
We will keep inventory in a concurrency-safe in-memory `Store`. All access goes through the `Store` methods
so persistence can be swapped in later without changing handlers.

## Options considered
1. **In-memory `Store`** - simple, zero dependencies; data lost on restart.
2. **SQLite** - durable; adds a dependency and migrations.
3. **Do nothing / mock only** - cannot demonstrate real behaviour.

## Consequences
- Positive: fast tests, no external services, standard library only.
- Negative: no durability; single-process only.
- Follow-up: introduce a `Repository` interface when persistence is required (new ADR).

## Notes for AI agents
- Never access the internal map outside `Store` methods.
- Do not introduce a database dependency without a new ADR.
