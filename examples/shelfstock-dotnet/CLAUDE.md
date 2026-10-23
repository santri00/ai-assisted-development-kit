# CLAUDE.md

@AGENTS.md

## Claude Code specifics
- Use plan mode for any change touching more than one project in the solution.
- Project slash commands live in `.claude/commands/` (`/spec`, `/review`).
- Hooks in `.claude/settings.json` run `task fmt` after edits; do not remove them.
