# PR Checklist (human + AI-assisted changes)

Use for every PR. For AI-assisted PRs the author is fully accountable: "the AI wrote it" is never a defence.

## A. Author (before requesting review)
- [ ] I read every line of the diff and can explain it.
- [ ] `task check` (build + lint + tests) passes locally.
- [ ] The PR does one thing; diff is < ~400 lines or split into stacked PRs.
- [ ] Linked spec / ticket; acceptance criteria are met.
- [ ] No unrelated refactors, reformatting, or drive-by changes slipped in by the agent.
- [ ] I ran the feature (or tests that exercise it), not just compiled it.

## B. Correctness
- [ ] Edge cases covered: empty, zero, negative, huge, duplicate, concurrent, unicode, timezone.
- [ ] Error paths handled and tested; no swallowed errors.
- [ ] New behaviour has tests; bug fix has a regression test that failed before the fix.
- [ ] Tests assert behaviour, not implementation details; no weakened/deleted assertions.
- [ ] No hallucinated APIs: every imported package/function actually exists and is used correctly.

## C. Design & maintainability
- [ ] Follows existing patterns and layering (see `docs/architecture.md`).
- [ ] No needless abstraction, dead code, or duplicated logic.
- [ ] Public API / schema / migration changes are backward compatible or called out.
- [ ] ADR added/updated for significant decisions.

## D. Security & compliance
- [ ] Inputs validated; outputs encoded; authn/authz checks present where needed.
- [ ] No secrets, tokens, keys, or PII in code, tests, fixtures, logs, or prompts.
- [ ] New dependencies reviewed: licence compatible, maintained, no known CVEs.
- [ ] No code copied verbatim from unknown sources (licence/IP risk).
- [ ] Logging does not leak sensitive data.

## E. Operability
- [ ] Observability: logs/metrics/traces for new failure modes.
- [ ] Config via env/secret manager; sensible defaults; documented.
- [ ] Rollout/rollback plan for risky changes (feature flag, migration order).

## F. Reviewer extras for AI-assisted PRs
- [ ] Looked for plausible-but-wrong logic (off-by-one, wrong units, swapped args).
- [ ] Checked tests were not written only to pass (tautologies, over-mocking).
- [ ] Confirmed agent did not modify CI, lint config, or test config without being asked.
- [ ] Asked "what would this agent have missed?" (business rules, hidden callers, data migrations).
