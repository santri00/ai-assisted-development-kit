# Prompt: Explore a codebase (read-only)

**When:** first time in a repo, or before touching an unfamiliar area.
**Mode:** read-only / plan mode. No edits.

```
Explore this repository and help me build an accurate mental model. Do NOT modify any files.

1. Read AGENTS.md, README, docs/architecture.md, and the build files.
2. Identify: entry points, main packages/modules and their responsibilities, data flow for the
   main use case, external dependencies/integrations, and how tests are organised.
3. Run the build and test commands from AGENTS.md and report the result.
4. Produce: (a) a 10-line architecture summary, (b) a list of 5 files I should read first,
   (c) risks or surprising things (dead code, missing tests, inconsistent patterns),
   (d) anything in AGENTS.md that appears wrong or outdated.

Cite file paths for every claim. If you are unsure, say so instead of guessing.
```
