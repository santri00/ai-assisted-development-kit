# Prompt: Implement one slice (test-first)

**When:** executing an approved plan, one step at a time.

```
Implement step <N> of the approved plan in docs/specs/<name>.md: "<STEP TITLE>".

Rules:
1. Write the failing test(s) first. Run them and show me they fail for the right reason.
2. Write the minimum code to make them pass. Follow existing patterns in neighbouring files.
3. Run `task check`. Fix failures properly - do not skip, delete, or weaken tests or linters.
4. Do not touch files outside the plan. If you need to, stop and ask.
5. Finish with: files changed, tests added, exact commands run + results, and anything uncertain.

Stop after this step. Do not start step <N+1>.
```
