# CLAUDE.md

@AGENTS.md

<!--
Claude Code reads CLAUDE.md (not AGENTS.md), so this file imports the shared instructions with
the @AGENTS.md line above. Add ONLY Claude-specific additions below. Keep the shared rules in AGENTS.md
so every tool gets the same guidance.
-->

## Claude Code specifics
- Use plan mode (Shift+Tab to cycle modes) for any change touching more than one package.
- Prefer the project slash commands in `.claude/commands/` (e.g. `/spec`, `/plan`, `/review`).
- Use subagents for read-only exploration and for independent review so the main context stays clean.
- Run `/compact` at natural breakpoints and `/clear` between unrelated tasks.
- Hooks in `.claude/settings.json` run the formatter after edits; do not remove them.
