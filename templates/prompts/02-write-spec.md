# Prompt: Write a spec (interview first)

**When:** before any feature larger than ~1 hour of work.
**Output:** `docs/specs/<name>.md` using `docs/specs/_template.md`.

```
I want to build: <ONE-PARAGRAPH FEATURE DESCRIPTION>

Before writing anything:
1. Read docs/architecture.md and the code this feature is likely to touch.
2. Interview me. Ask up to 8 clarifying questions, most important first, one batch at a time.
   Cover: users/use cases, inputs/outputs, edge cases, failure behaviour, performance, security,
   backward compatibility, and what is explicitly out of scope.
3. After I answer, write the spec using docs/specs/_template.md.
   Acceptance criteria must be testable (Given/When/Then or concrete examples).
4. List open questions and assumptions at the bottom. Do not write implementation code.
```
