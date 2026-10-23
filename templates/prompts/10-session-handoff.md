# Prompt: Session handoff (context reset without losing progress)

**When:** context is getting long, switching tools (Claude <-> Codex), or ending the day.

```
Write a handoff note to docs/handoff/<date>-<topic>.md (create the folder if needed) so a fresh
agent session can continue without me re-explaining.

Include:
1. Goal and link to the spec/plan.
2. Done so far (with commit hashes or file paths).
3. Current state: does `task check` pass? Paste the last result.
4. Next 1-3 concrete steps.
5. Decisions made and why; approaches tried that failed.
6. Open questions / risks / things not to touch.

Be factual and terse. Do not include secrets or customer data.
```
