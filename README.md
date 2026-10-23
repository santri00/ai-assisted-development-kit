# AI-Assisted Development Kit

A tool-agnostic, language-agnostic kit for adopting AI coding agents (Claude Code, Codex, Copilot) in an engineering team. Two runnable reference repos (Go and .NET), a Taskfile-based command vocabulary, and Docker/Docker Desktop support.

| Start here | Purpose |
|---|---|
| [`HANDBOOK.md`](HANDBOOK.md) | The full step-by-step learning guide (Steps 0-13, plus 4b Containers) with worked examples, exercises, and an honest Taskfile evaluation |
| [`PLAYBOOK.md`](PLAYBOOK.md) | The 10-minute team reference: rules, workflow, definition of done, security non-negotiables, Docker Desktop cheat sheet |
| [`templates/`](templates/) | Copy-paste starters: `AGENTS.md`, `CLAUDE.md`, `Taskfile.yaml`, tool configs, Docker files, PR checklist, prompt library, ADR/spec/architecture templates, security policy |
| [`examples/shelfstock-go/`](examples/shelfstock-go/) | Go reference service with the finished scaffold (Taskfile, Dockerfile, compose, CI) |
| [`examples/shelfstock-dotnet/`](examples/shelfstock-dotnet/) | .NET (ASP.NET Core minimal API) reference service with the same scaffold |
| `examples/low-stock-reference-solution-{go,dotnet}.patch` | Answer keys for the practice feature |

## Layout

```
ai-assisted-development-kit/
├── README.md  HANDBOOK.md  PLAYBOOK.md
├── templates/
│   ├── AGENTS.md  CLAUDE.md  PR-CHECKLIST.md  Taskfile.yaml  .gitattributes.template
│   ├── claude/    settings.json, commands/{spec,review}.md
│   ├── codex/     config.toml
│   ├── github/    pull_request_template.md, copilot-instructions.md, ci.yml
│   ├── docker/    Dockerfile.go, Dockerfile.dotnet, compose.yaml, devcontainer.json, dockerignore.template
│   ├── dotnet/    .editorconfig, Directory.Build.props, global.json
│   ├── prompts/   01-explore ... 10-session-handoff
│   └── docs/      architecture.md, adr-template.md, spec-template.md, ai-security-and-compliance-policy.md
└── examples/
    ├── shelfstock-go/       cmd/, internal/, Taskfile.yaml, Dockerfile, compose.yaml, AGENTS.md, docs/ ...
    ├── shelfstock-dotnet/   src/{Core,Api}, tests/, Taskfile.yaml, Dockerfile, compose.yaml, AGENTS.md, docs/ ...
    └── low-stock-reference-solution-{go,dotnet}.patch
```

## 5-minute start

```bash
winget install Task.Task        # or: brew install go-task   (Task replaces make; works on Windows)

cd examples/shelfstock-go       # Go track
task check                      # vet + gofmt + tests + build

cd ../shelfstock-dotnet         # .NET track (needs .NET SDK 10; see HANDBOOK for other SDKs)
task fmt && task check

task docker:check               # same checks inside Linux Docker (needs Docker Desktop running)
task docker:up                  # http://localhost:8080/healthz
```

## Verification status (be aware)

| Item | Status |
|---|---|
| Go example, Taskfiles, Go reference patch | Executed and passing while authoring |
| **.NET example** | **Not compiled or tested** (no .NET SDK / NuGet in the authoring environment). Expect small first-run fixes (SDK version, package versions, `dotnet format` spacing) |
| Dockerfiles / compose | Parsed and partially exercised; **images were not built** (no Docker daemon while authoring) |

## Adopting in a real repository

1. Copy `templates/AGENTS.md` to the repo root and fill it in (keep it short).
2. Copy `templates/Taskfile.yaml` and fill in the task bodies; keep the task names.
3. Copy `templates/CLAUDE.md`, `templates/claude/settings.json` -> `.claude/`, `templates/codex/config.toml` -> `.codex/`.
4. Pick `templates/docker/Dockerfile.go` or `.dotnet`, plus `compose.yaml`; add `.gitattributes` from the template.
5. Wire CI via `templates/github/ci.yml`; add `templates/docs/*`, `PR-CHECKLIST.md`, and the PR template.
6. Review `templates/docs/ai-security-and-compliance-policy.md` (including the Docker section) with Security and Legal.

Tool-specific settings change quickly; confirm them against the vendor docs linked in HANDBOOK Appendix C.
