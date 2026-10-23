# CLAUDE.md

@AGENTS.md

## Claude Code specifics
- Use plan mode for any change touching more than one package.
- Project slash commands live in `.claude/commands/` (`/spec`, `/review`).
- Hooks in `.claude/settings.json` run `task fmt` after edits; do not remove them.
- On Windows without a C compiler, `task test` omits `-race`; run `task docker:check` for the full race run.
