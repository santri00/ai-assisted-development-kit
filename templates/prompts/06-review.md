# Prompt: Independent code review

**When:** before opening a PR, and as an automated CI step. Use a *fresh* session (no authoring context).

```
Review the diff of this branch against main (git diff main...HEAD). Do NOT edit files.

Check, in priority order:
1. Correctness: logic errors, edge cases, concurrency, error handling, resource leaks.
2. Security: injection, authz gaps, secrets, unsafe deserialization, SSRF, path traversal, PII in logs.
3. Spec conformance: compare against docs/specs/<name>.md acceptance criteria.
4. Tests: missing cases, tautological assertions, weakened or deleted tests.
5. Maintainability: unnecessary complexity, pattern drift, dead code, API compatibility.
6. Hallucination check: verify every external API/function/flag used actually exists in this repo's versions.

Output a table: severity (Blocker/Should fix/Nit) | file:line | issue | suggested fix.
Be specific and skeptical. If the diff is fine, say what you verified rather than inventing issues.
```
