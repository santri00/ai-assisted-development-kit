# Prompt: Safe refactor

**When:** improving structure without changing behaviour.

```
Refactor <TARGET> to <GOAL, e.g. "separate persistence from business logic">.

Constraints:
1. Behaviour must not change. First confirm test coverage is adequate; if not, stop and propose
   characterization tests (prompt 05) before refactoring.
2. Work in small steps; after each step run `task check` and keep the repo green.
3. No public API changes unless I approve. No new dependencies.
4. Commit-sized steps with a one-line description each, so I can review them individually.
5. At the end summarise: before/after structure, files moved/renamed, and how behaviour
   equivalence was verified.
```
