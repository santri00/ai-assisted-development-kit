# The AI-Assisted Development Handbook (Example)

A step-by-step, tool-agnostic, language-agnostic guide to adopting AI coding agents (Claude Code, Codex, Copilot and others) in a professional engineering team.

- **Audience:** an engineering team.
- **Running examples (pick a track, or do both):** `examples/shelfstock-go` (Go) and `examples/shelfstock-dotnet` (C# / ASP.NET Core). Same service, same spec, same scaffold, so you can compare how the method carries across languages.
- **Command runner:** [Task](https://taskfile.dev) (`Taskfile.yaml`) in both examples, with identical task names. Step 1 evaluates it honestly, including its shortcomings.
- **Local runtime:** Dockerfile + Docker Compose, designed for teams on Docker Desktop (including Windows). See Step 4b.
- **Companion files:** `PLAYBOOK.md` (the short team reference) and `templates/` (copy-paste starters).
- **Tool focus:** Claude Code and Codex. Copilot is covered lightly. The scaffolding works for all three.

> **A note on freshness.** Tool details (flags, config keys, slash commands) change quickly. The *method* in this handbook is stable; wherever a specific tool setting appears, confirm it in the vendor docs linked in [Appendix C](#appendix-c--further-reading) before rolling it out.

---

## Table of contents

1. [How to use this handbook](#1-how-to-use-this-handbook)
2. [Foundations: the mental model](#2-foundations-the-mental-model)
3. [The core idea: one standard scaffold](#3-the-core-idea-one-standard-scaffold)
4. [The steps](#4-the-steps)
   - [Step 0: Set up your tools safely](#step-0--set-up-your-tools-safely)
   - [Step 1: Make the repo agent-ready (one command to prove "done")](#step-1--make-the-repo-agent-ready)
   - [Step 2: Write `AGENTS.md`](#step-2--write-agentsmd)
   - [Step 3: Give the agent architecture and decision context](#step-3--give-the-agent-architecture-and-decision-context)
   - [Step 4: Add guardrails (permissions, sandbox, hooks)](#step-4--add-guardrails)
   - [Step 4b: Containers and Docker Desktop](#step-4b--containers-and-docker-desktop)
   - [Step 5: Explore before you change anything](#step-5--explore-before-you-change-anything)
   - [Step 6: Spec, plan, implement (the core workflow)](#step-6--spec-plan-implement)
   - [Step 7: Test-first slices and verification loops](#step-7--test-first-slices-and-verification-loops)
   - [Step 8: Review and quality gates](#step-8--review-and-quality-gates)
   - [Step 9: Build reusable assets (prompts, commands, skills, subagents, MCP)](#step-9--build-reusable-assets)
   - [Step 10: Manage context and long-running work](#step-10--manage-context-and-long-running-work)
   - [Step 11: Security, IP and compliance](#step-11--security-ip-and-compliance)
   - [Step 12: Measure what matters](#step-12--measure-what-matters)
   - [Step 13: Roll out to a team](#step-13--roll-out-to-a-team)
5. [Cross-tool mapping (Claude Code / Codex / Copilot)](#5-cross-tool-mapping)
6. [Porting to other languages](#6-porting-to-other-languages)
7. [Anti-patterns and troubleshooting](#7-anti-patterns-and-troubleshooting)
8. [Appendices](#appendices)

---

## 1. How to use this handbook

**Suggested pace: about two weeks of part-time practice.**

| Days | Steps | Outcome |
|---|---|---|
| 1 | 0-2 | A repo an agent can build, test, and understand |
| 2 | 3-4 | Context docs and safe permissions |
| 3-5 | 5-7 | You complete one real feature with the spec/plan/implement loop |
| 6-7 | 8-9 | Review gates and reusable prompts/commands |
| 8-9 | 10-12 | Long-running work, security posture, metrics |
| 10 | 13 | A draft rollout plan for your team |

**How each step is laid out:** Goal, Why it matters, What to do, Worked example (shelfstock), Exercise, Done when.

**Prerequisites**

| Tool | Needed for | Install |
|---|---|---|
| Git | everything | https://git-scm.com |
| [Task](https://taskfile.dev) 3.x | `task check` etc. (replaces `make`; works on Windows) | `winget install Task.Task`, `brew install go-task`, or `go install github.com/go-task/task/v3/cmd/task@latest` |
| Go 1.22+ | Go track | https://go.dev/dl |
| .NET SDK 10 | .NET track (change `global.json` + `Directory.Build.props` if you are on another SDK) | https://dot.net |
| Docker Desktop (or Docker Engine) | Step 4b, `task docker:*` | https://docs.docker.com/desktop |
| Claude Code and/or Codex | the agents | see Step 0 |

**Verification status of this kit (read this once)**

| Component | Status |
|---|---|
| Go example: `task check`, failing-lint negative test, Windows-style no-race path, reference patch applied to a clean copy | Executed and passing (Go 1.24, Task 3.54) |
| Taskfile template and all YAML/JSON/XML files | Parsed successfully; the Go and template Taskfiles also run under Task |
| Go `Dockerfile` / `compose.yaml` | Compose file parses; the static binary builds with the Dockerfile's exact flags. **Image not built** |
| **.NET example (all of it)** | **Not compiled or tested** |
| .NET `Dockerfile` / `compose.yaml` | **Image not built** |

On your machine, the first run of the .NET track is: `task fmt`, then `task check`. Likeliest first-run adjustments: SDK/TFM mismatch (`global.json`, `Directory.Build.props`), a package version, or `dotnet format` re-spacing. Fixing those is itself a good exercise in using an agent with a failing check as the steering wheel.

**Setting up the practice repo** (choose one track)

```bash
# Go track
cd examples/shelfstock-go
task check          # vet + gofmt check + tests (-race where supported) + build. Must pass before you start.

# .NET track
cd examples/shelfstock-dotnet
task fmt            # once, to settle formatting
task check          # format check + build (warnings are errors) + tests

git init && git add -A && git commit -m "baseline"
```

Each example already contains the *finished* scaffold (AGENTS.md, docs, settings) so you can read what "good" looks like. For hands-on learning, **delete the scaffold files and rebuild them yourself step by step**, using the originals as an answer key:

```bash
# Optional: strip the scaffold to practise building it
rm -rf AGENTS.md CLAUDE.md .claude .codex .github docs      # keep Taskfile.yaml, Dockerfile, compose.yaml: those are Steps 1 and 4b
```

The practice feature is **low-stock alerts** (`docs/specs/low-stock-alerts.md`). A reference solution for each track is provided as a patch (`examples/low-stock-reference-solution-go.patch`, `examples/low-stock-reference-solution-dotnet.patch`; the .NET one is uncompiled). Do not peek until you have tried it.

---

## 2. Foundations: the mental model

### 2.1 What a coding agent actually is

An agent is **a language model + tools + a loop**:

1. It reads context (your prompt, instruction files, files it opens).
2. It decides on an action (read a file, edit, run a command, search).
3. A tool executes the action and returns output.
4. The output goes back into context, and the loop repeats until it decides it is done.

Everything useful follows from this:

| Fact | Consequence for you |
|---|---|
| It only knows what is in its context window | Put durable knowledge in files it will read (`AGENTS.md`, docs), not in your head |
| Context fills up and quality degrades | Keep tasks small; reset between tasks; hand off with notes |
| It learns about success from tool output | Give it fast, trustworthy checks (`task check`). **The check is the steering wheel.** |
| It can run commands and edit files | You need permissions, sandboxes, and review, the same as for any junior with root |
| It is fluent and confident even when wrong | Verify; never accept "done" without evidence |
| It follows instructions found anywhere in context | Untrusted content (web pages, issues) can hijack it (prompt injection) |

### 2.2 The maturity ladder

| Level | Description | Typical tool use | Risk |
|---|---|---|---|
| 1. Autocomplete | Inline suggestions | Copilot completions | Low |
| 2. Chat | Ask questions, paste snippets | Chat panes | Low-Med |
| 3. Agent, supervised | Agent edits and runs commands; human approves and reviews | Claude Code, Codex CLI | Med |
| 4. Agent, directed with a harness | Spec-driven; agent verifies itself against tests/CI; human reviews outcomes | Same tools + standard scaffold | Med |
| 5. Agent, automated | Agents triggered by tickets/CI; produce PRs unattended | Cloud agents, GitHub/Bitbucket integration, headless CLI | Med-High |

**This handbook takes you to Level 4 and shows how to approach Level 5 safely.** Skipping to Level 5 without the scaffolding of Levels 3-4 is the most common reason rollouts fail.

### 2.3 What agents are good and bad at

| Strong | Weak / risky |
|---|---|
| Boilerplate, CRUD, glue code | Subtle business rules not written down anywhere |
| Writing and extending tests | Large cross-cutting redesigns in one shot |
| Codebase exploration and summarising | Performance tuning without benchmarks |
| Mechanical refactors and migrations | Security-critical logic (auth, crypto) without expert review |
| First drafts of docs, specs, ADRs | Knowing what is *not* in the repo (production behaviour, tribal knowledge) |
| Reproducing and bisecting bugs | Judging product trade-offs |

### 2.4 The four principles used throughout

1. **Context is the product.** Better written context beats better prompting.
2. **Small verified steps.** Every step ends with a passing check.
3. **Humans own decisions and accountability.** Agents propose; people decide and sign.
4. **Standardize the scaffold, not the tool.** Files and commands that any agent can use.

---

## 3. The core idea: one standard scaffold

The same repository layout should work with any agent. Only a thin tool-specific layer differs.

```
repo/
├── AGENTS.md                     <- THE source of truth for all agents
├── CLAUDE.md                     <- one line: @AGENTS.md  (+ Claude-only notes)
├── Taskfile.yaml                 <- one stable command vocabulary: build, test, lint, fmt, check, run, docker:*
│                                    (or Makefile / package scripts: the NAMES matter more than the runner)
├── Dockerfile, compose.yaml      <- reproducible build/test/run; the `test` stage mirrors CI
├── .gitattributes, .editorconfig <- LF line endings + formatting rules (essential with Windows developers)
├── docs/
│   ├── architecture.md           <- 1-2 page system map
│   ├── adr/                      <- decision records
│   ├── specs/                    <- feature specs (input to agents)
│   └── PR-CHECKLIST.md
├── .github/
│   ├── pull_request_template.md
│   ├── copilot-instructions.md   <- pointer to AGENTS.md
│   └── workflows/ci.yml          <- runs the same `task check`
├── .claude/                      <- Claude Code only: settings.json, commands/, skills/, agents/
└── .codex/                       <- Codex only: config.toml
```

**The four portable pillars**

| Pillar | Artifact | Why it is portable |
|---|---|---|
| Instructions | `AGENTS.md` | Plain Markdown. Read natively by Codex and other agents; Claude Code imports it with `@AGENTS.md`; Copilot can read it or be pointed to it |
| Verification | `task check` | Any agent can run a shell command |
| Environment | `Dockerfile` (`test` + `runtime` stages) and `compose.yaml` | The same Linux build on Windows, macOS, CI and cloud agents; Docker Desktop runs it locally |
| Process | Spec -> plan -> small slices -> review | It is a way of working, not a feature of a product |

The tool-specific files (`.claude/`, `.codex/`) hold **permissions and conveniences only**, never the rules of the project. If you moved to a new tool tomorrow, you would lose only convenience.

---

## 4. The steps

### Step 0 · Set up your tools safely

**Goal:** have Claude Code and Codex installed, authenticated with the account, and running in restricted modes.

**Why it matters:** Agents run commands on your machine. Start restrictive; loosen deliberately.

**What to do**

1. Install and authenticate through approved route (SSO / enterprise plan). Do not use personal accounts on company code. See [Step 11](#step-11--security-ip-and-compliance).
2. Learn the three permission postures every tool offers in some form:

| Posture | Meaning | Claude Code | Codex |
|---|---|---|---|
| Read-only / plan | Can read and propose; cannot change anything | Plan mode (`Shift+Tab` cycles modes, or `--permission-mode plan`) | `sandbox_mode = "read-only"` |
| Supervised edits | Edits allowed; commands need approval | Default mode with allow/deny rules | `approval_policy = "on-request"` + `workspace-write` |
| Full access | No prompts | "Bypass permissions" | `danger-full-access` / `approval_policy = "never"` |

3. **Rule:** never use full access outside a disposable container or VM with no production credentials.
4. Learn the five commands you will use daily in your tool: start a session, switch to plan mode, compact/clear context, run a one-off non-interactive prompt (`claude -p "..."`, `codex exec "..."`), and view/resume sessions.

**Worked example:** start Claude Code inside `examples/shelfstock-go` in plan mode and ask: "What does this repo do? Do not edit anything." Then do the same in Codex with the read-only sandbox. Compare answers. They will differ, and that difference is why shared written context (Steps 2-3) matters.

**Exercise:** try to get the agent to edit a file while in read-only/plan mode. Observe how it refuses or asks. Now you know what the guardrail feels like.

**Done when:** you can start each tool in read-only mode and in supervised-edit mode on demand.

---

### Step 1 · Make the repo agent-ready

**Goal:** one command that tells anyone (human, CI, agent) whether the repo is healthy.

**Why it matters:** Agents improve by running checks and reading failures. If "build and test" takes six undocumented commands, or the suite is flaky or takes 30 minutes, the agent will guess, skip, or declare victory falsely. **This is the highest-leverage step in the handbook.**

**What to do**

1. Create a stable command vocabulary. Whatever your language, expose the same names:

| Target | Meaning |
|---|---|
| `build` | compile / package |
| `test` | fast tests (aim for < 2 minutes) |
| `lint` | static checks, formatting verification |
| `fmt` | auto-format |
| `check` | `lint` + `test` + `build`: **the definition of done** |
| `run` | start locally |

2. Make `check` fast, deterministic, and hermetic (no network, no shared state, no flaky tests). Fix or quarantine flaky tests now.
3. Run it in CI exactly as-is (Step 8).
4. Add a README section "Getting started" that works from a clean clone.

**Worked example, Go (`examples/shelfstock-go/Taskfile.yaml`):**

```yaml
version: '3'
vars:
  # race detector needs cgo + a C compiler: fine on Linux/macOS/CI/Docker, usually absent on Windows
  RACE: '{{if eq OS "windows"}}{{else}}-race{{end}}'
tasks:
  check:
    desc: "FULL PRE-MERGE CHECK (definition of done): lint + test + build"
    cmds: [ {task: lint}, {task: test}, {task: build} ]
  test:
    cmds: [ "go test {{.RACE}} ./..." ]
  lint:
    cmds:
      - go vet ./...
      - |
        out="$(gofmt -l .)"
        if [ -n "$out" ]; then echo "gofmt needed:"; echo "$out"; exit 1; fi
```

Run `task check`. The `-race` flag matters: `Store` is concurrent, and an agent that "fixes" locking incorrectly gets caught here. On Windows without a C compiler the task drops `-race`; run `task docker:check` (Linux container) to get the race run anyway.

**Worked example, .NET (`examples/shelfstock-dotnet/Taskfile.yaml`):** same task names, different bodies:

```yaml
tasks:
  lint:  { cmds: [ "dotnet format --verify-no-changes" ] }
  build: { deps: [restore], cmds: [ "dotnet build --no-restore --nologo" ] }   # warnings are errors
  test:  { deps: [restore], cmds: [ "dotnet test --no-restore --nologo" ] }
  check: { cmds: [ {task: lint}, {task: build}, {task: test} ] }
```

`Directory.Build.props` turns on nullable reference types and `TreatWarningsAsErrors`, so the compiler itself becomes part of the steering wheel: an agent cannot "fix" a nullability warning by ignoring it.

#### 1.1 Choosing the command runner: Task vs Make vs the rest

The *names* (`check`, `test`, ...) are the standard. The runner is an implementation detail you can swap. This kit uses Task as an example because I am using Windows environment with the Docker Desktop, and `make` is not available on Windows by default.

| Option | Strengths | Weaknesses | Good fit |
|---|---|---|---|
| **Task** (`Taskfile.yaml`) | Single binary; works natively on Windows/macOS/Linux; YAML; `desc` makes `task --list` self-documenting (useful to agents); `deps`, vars, `includes` for monorepos; built-in portable shell | Extra install; YAML quoting traps; not real bash (see below) | Polyglot or Windows-heavy teams (this kit) |
| **Make** | Preinstalled on Linux/macOS; universally understood by humans and models | Not on Windows by default (needs WSL/Git Bash/choco); tab-sensitive; `$$` escaping; odd semantics for non-file tasks | Linux/macOS-only teams |
| **just** | Make-like but simpler, cross-platform | Another binary; smaller ecosystem | Individuals who like Make syntax |
| **mise tasks** | Pins tool versions AND runs tasks | Larger concept surface | Teams wanting version pinning built in |
| **Native** (`npm run`, `dotnet`, `gradlew`) | Zero extra tooling | Names differ per ecosystem; no uniform vocabulary across repos | Single-language orgs; also valid for .NET (`task` here is only a thin wrapper) |
| **Scripts** (`.sh` + `.ps1`) | No dependency | Two scripts to keep in sync; drift | Last resort |

**Shortcomings of Taskfile.yaml to plan for**:

| # | Shortcoming | Why it matters | Mitigation |
|---|---|---|---|
| 1 | **Another prerequisite** that must exist on every laptop, CI image, dev container and agent sandbox | A cloud agent sandbox with restricted network may be unable to install it; `make` is already there on Linux | Install in devcontainer/CI (`arduino/setup-task`); document install in `AGENTS.md`; list the underlying raw commands as a fallback |
| 2 | **YAML traps.** An unquoted `echo "TODO: x"` fails with `invalid keys in command` (it parses `TODO: x` as a map). Multi-line `|` scripts echo noisily | Agents edit these files; a quoting slip breaks the whole vocabulary | Quote command strings containing `: `; run `task --list` in CI as a parse smoke test; keep tasks thin and move complex logic into scripts |
| 3 | **Embedded shell is POSIX-like, not bash** | No bash-only features; behaviour can differ from your terminal; the tools you call (`gofmt`, `dotnet`, `docker`) must still be on the Windows PATH | Keep scripts simple; test on Windows once |
| 4 | **Name collision**: Taskwarrior also installs a `task` binary on some Linux distros | `task check` may run the wrong program | Use the `go-task` package/alias; assert `task --version` in CI |
| 5 | **Version drift**: the schema is `version: '3'` but features vary by minor release | Different laptops, different behaviour | Pin the exact version in CI/devcontainer and state the minimum in the README |
| 6 | **Per-OS conditionals are clunky** (`{{if eq OS "windows"}}`) | Platform differences leak into YAML (we needed one for Go's `-race`) | Prefer containers for parity (`task docker:check`) over OS branches |
| 7 | **`deps` run in parallel** | Surprising ordering, interleaved output | Use sequential `cmds: - task: x` when order matters (as `check` does) |
| 8 | **Less universal than `make`/`npm`** for models and engineers | Agents sometimes reach for `make` or raw commands | The command table in `AGENTS.md`; allow-list only `task` targets |
| 9 | **Indirection hides what runs. Security-relevant.** An allow rule such as `Bash(task docker:check)` permits *whatever that target does*, and agents can edit `Taskfile.yaml` | A deny rule on `dotnet nuget push` is bypassed if a task wraps it; an agent could add a malicious command to a task | Treat `Taskfile.yaml`, `Dockerfile`, `compose.yaml` and CI as **security-sensitive**: CODEOWNERS review, "Ask first" in `AGENTS.md`, never auto-approve edits to them |
| 10 | **Overkill for single-language repos** where one native command already does the job | An extra moving part for little gain | Fine to call `dotnet`/`npm` directly if the team never needs a uniform vocabulary; keep the *names* anyway |

**Verdict:** keep Task for this kit (cross-platform, self-documenting, uniform across Go and .NET), keep the Taskfiles thin wrappers over native commands, pin the version, and protect the file with review rules. The vocabulary is the asset; the runner is replaceable.

**Exercise:** introduce a deliberate failure (change an expected value in a test). Confirm `task check` fails clearly. Revert. A check that cannot fail is useless.

**Done when:** run `task check`, and trust the result.

---

### Step 2 · Write `AGENTS.md`

**Goal:** a short, accurate instruction file that every agent reads at the start of every session.

**Why it matters:** This is the agent's onboarding doc. Agents have no memory between sessions; this file is their memory of your project.

**The design rules**

| Rule | Reason |
|---|---|
| **Short** (aim < 150 lines; Codex truncates large instruction files by default at ~32 KiB) | Long files get ignored or truncated; every line costs context |
| **Only what the agent would get wrong** | Do not document what the code already shows |
| **Commands, not prose** | Executable beats descriptive |
| **Boundaries as Always / Ask first / Never** | Clear enforcement tiers |
| **Living document** | Every time an agent makes a repeatable mistake, add one line |
| **No secrets, ever** | The file is committed |

**What to do**

1. Start from `templates/AGENTS.md`. (Both Claude Code's and Codex's `/init` commands can generate a draft from your repo. Treat that as a first draft and prune it hard.)
2. Fill in: overview, commands, repo map, conventions, boundaries, workflow, security, gotchas.
3. Create `CLAUDE.md` containing `@AGENTS.md` so Claude Code reads the same rules. Create `.github/copilot-instructions.md` pointing at `AGENTS.md` for Copilot.
4. For monorepos, add nested `AGENTS.md` files in subdirectories with local rules (most specific wins). Keep the root one generic.

**Worked example:** read `examples/shelfstock-go/AGENTS.md` and compare it with `examples/shelfstock-dotnet/AGENTS.md`. Notice (Go version):

- **Layering rule** ("`inventory` must never import `net/http`"): an architectural rule an agent cannot infer from one file. In .NET the same rule is enforced by the compiler through project references (`ShelfStock.Core` has none), so the `AGENTS.md` only needs to say "never add a reference from Core to Api".
- **Gotcha** ("`Remove` on missing SKU is 404, but insufficient stock is 409"): a distinction agents commonly collapse.
- **Definition of done**: "`task check` passes. Never claim work is complete without running it."

Now test it. In a fresh session, with no other context, ask the agent: "Add a `GET /stats` endpoint returning the number of SKUs and total quantity." Check that it (a) kept domain logic in the core (`internal/inventory` in Go, `ShelfStock.Core` in .NET) and only HTTP concerns in the transport layer (`httpapi` / `ShelfStock.Api`), (b) wrote tests in the project's style (table-driven in Go, `[Theory]` in xUnit), (c) ran `task check`. If not, the file is missing something. Add a line and retry.

**Exercise:** write your `AGENTS.md` for shelfstock from the template without looking at the finished one. Then diff against the answer key.

**Done when:** a fresh agent session follows your conventions on a small task without being reminded.

---

### Step 3 · Give the agent architecture and decision context

**Goal:** a short architecture doc and decision records the agent can read.

**Why it matters:** `AGENTS.md` says *how to work*; architecture docs say *what the system is*. Without them the agent reinvents patterns, violates boundaries, and re-litigates settled decisions.

**What to do**

1. Write `docs/architecture.md` (template provided): purpose, components and responsibilities, 1-3 key flows, data, cross-cutting rules, known debt. Aim for 1-2 pages. Include a diagram in Mermaid text form so agents can read it (use C4 levels 1-2 if your team uses C4).
2. Start `docs/adr/` with the ADR template. Record decisions that are easy to violate or reverse. Add a **"Notes for AI agents"** section to each ADR with enforceable one-line rules.
3. Link both from `AGENTS.md`. Do not duplicate content; link.
4. Add the rule: "Update architecture docs and ADRs in the same PR that changes them."

**Worked example:** `docs/adr/0001-in-memory-store.md` records *why* storage is in-memory and gives the agent an enforceable rule: "Do not introduce a database dependency without a new ADR." Ask an agent to "persist stock to SQLite." A well-scaffolded agent should stop and propose an ADR first rather than silently adding a dependency.

**Exercise:** write an ADR for a decision you made recently in a real project. Add the "Notes for AI agents" section. Which rule would an agent most likely break without it?

**Done when:** an agent asked to make a structural change cites the architecture doc or ADR in its plan.

---

### Step 4 · Add guardrails

**Goal:** make the safe path the default, and make dangerous actions impossible or require approval.

**Why it matters:** Instructions are advisory; permissions are enforced. Use both.

**Layers of defence**

| Layer | Mechanism | Example |
|---|---|---|
| Instructions | `AGENTS.md` "Never" list | "Never commit secrets" |
| Tool permissions | `.claude/settings.json` allow/deny; Codex `approval_policy` | Deny `Read(./.env)`, deny `git push` |
| Sandbox | Codex `sandbox_mode`; containers/devboxes | No network, write only in workspace |
| Hooks (deterministic) | Claude Code hooks run commands on events | Auto-format after every edit; block edits to protected paths |
| Platform | Branch protection, CODEOWNERS, required checks, secret scanning | Agents cannot merge without review |
| Credentials | Least-privilege, short-lived tokens | Agents never get production credentials |
| Containers | Docker/dev container boundaries, allow-listed `task docker:*` targets, no Docker socket for agents | Agent in a disposable dev container; raw `docker run` denied (Step 4b) |

**What to do**

1. Copy `templates/claude/settings.json` to `.claude/settings.json`. Customise allow rules to your `task` targets (and remember that allowing `task x` allows everything `x` does: see shortcoming 9 in Step 1.1); keep `deny` for `.env`, secrets, `rm -rf`, `git push`, network fetchers.
2. Copy `templates/codex/config.toml` to `.codex/config.toml` (or your user config) and confirm the keys against the current docs. The baseline: `approval_policy = "on-request"`, `sandbox_mode = "workspace-write"`, network off.
3. Commit `.claude/settings.json` (shared team rules). Keep personal overrides in `.claude/settings.local.json` and git-ignore it.
4. Turn on branch protection: required PR review, required `check` status, no force-push to main.
5. Use a secret scanner locally and in CI, and enable push protection if your host supports it.

**Hooks are the underrated feature.** Prompts are suggestions; hooks always run. The shelfstock example formats Go code after every edit:

```json
"hooks": { "PostToolUse": [ { "matcher": "Edit|Write",
  "hooks": [ { "type": "command", "command": "task fmt >/dev/null 2>&1 || true" } ] } ] }
```

**Exercise:** ask the agent to "read .env and summarise it." Confirm the deny rule blocks it. Then ask it to run `git push`. Confirm it is blocked or requires approval.

**Done when:** the dangerous actions on your list fail even if you *ask* the agent to do them.

---

### Step 4b · Containers and Docker Desktop

**Goal:** make local runs and checks identical for everyone (Windows, macOS, Linux, CI), and decide deliberately how the agent relates to Docker.

**Why it matters:** "works on my machine" is worse with agents: they run commands on whatever machine they are on. A container gives you one reproducible Linux environment, and Docker access is also one of the most dangerous permissions you can hand an agent.

**Three operating models**

| Model | Where the agent runs | Docker's role | Pros | Cons | Use when |
|---|---|---|---|---|---|
| **A. Host-native + Docker for parity** (start here) | On the developer's machine | Docker Desktop runs the app (`task docker:up`) and parity checks (`task docker:check`); agent only gets allow-listed `task docker:*` targets | Simple; IDE-friendly; native speed | Agent has host access; toolchain drift between laptops | Phases 1-3 of the rollout |
| **B. Agent inside a dev container / VM** | In a container with only the repo mounted (`templates/docker/devcontainer.json`) | The container *is* the toolchain and the sandbox; Docker-outside-of-Docker only if essential | Hermetic; same Linux everywhere; small blast radius | Setup effort; slower bind mounts on Windows unless the repo lives in the WSL2 filesystem; credentials must be passed in deliberately | Higher-autonomy work; any use of full-access modes |
| **C. Cloud sandbox / CI agent** | Hosted, isolated environment | CI builds images; sandbox may have no Docker daemon | Isolation and parallelism | Network allow-lists and secrets design needed | Phase 4 automation |

> Docker is also shipping features aimed at running coding agents in isolated sandboxes. Evaluate them against Model B, and verify current availability and maturity in the Docker docs before depending on them.

**The pattern: Docker as the contract.** Both example repos have a multi-stage `Dockerfile` with two targets that matter:

| Target | What it does | Used by |
|---|---|---|
| `test` | Runs the checks inside Linux (`go vet` + `go test -race`, or `dotnet format --verify-no-changes` + `dotnet test`). **The build fails if checks fail** | `task docker:check`, CI |
| `runtime` | Small, non-root production image | `task docker:build`, `compose.yaml`, deployment |

`task docker:check` therefore gives a Windows developer the same result as CI, including the Go race detector that needs a C compiler on the host.

**What to do**

1. Copy `templates/docker/Dockerfile.go` or `Dockerfile.dotnet`, `compose.yaml`, and `dockerignore.template` (as `.dockerignore`). Keep the `test` stage.
2. Add the `docker:*` tasks (see `templates/Taskfile.yaml`): `docker:check`, `docker:build`, `docker:up`, `docker:logs`, `docker:down` (and `docker:watch` for Compose Watch rebuilds).
3. Add `.gitattributes` (`templates/.gitattributes.template`). **Line endings are the number one Windows trap**: CRLF in a `Dockerfile` or shell script breaks it inside a Linux container, and `dotnet format --verify-no-changes` fails on CRLF files.
4. Give the agent **narrow, project-defined entry points**, not raw Docker:
   - Allow: `task docker:check` (and optionally `task docker:build`).
   - Deny: `docker run`, `docker system prune`, `docker volume rm`, `docker compose down -v` (already in `templates/claude/settings.json`).
   - Never mount the Docker socket (`/var/run/docker.sock`, or the Docker Desktop pipe on Windows) into a container an agent controls: it is equivalent to root on the host.
5. Keep `compose.yaml` safe by default: `read_only: true` (+ `tmpfs: /tmp` for .NET), `no-new-privileges`, non-root user in the image, secrets from an untracked `.env` or a secret manager, never inline.
6. Pin base image tags (digests for production), scan images in CI, and add `docker build --target test .` as a CI job (the example `ci.yml` files do).

**Windows + Docker Desktop checklist**

- Use the **WSL2 backend**. For heavy bind-mount workloads keep the repo inside the WSL filesystem (`\\wsl$\...`) rather than `C:\`: file I/O across the boundary is much slower.
- **Run the agent and the repo in the same environment** (both in WSL, or both on Windows). Mixing causes path-translation errors. Check each vendor's current guidance on native Windows versus WSL (at the time of writing Claude Code supports both; Codex's Windows support has been experimental with WSL recommended: confirm).
- **No `make` on Windows**: that is why this kit uses Task (`winget install Task.Task`).
- **`.gitattributes` with `eol=lf`** and `.editorconfig` with `end_of_line = lf`.
- Task runs commands in its own portable shell, so `$(...)` in a Taskfile works from PowerShell; the tools you call (`go`, `dotnet`, `docker`) must be on the PATH.
- Go's `-race` needs cgo and gcc on Windows: use `task docker:check` instead of installing a C toolchain.
- Docker Desktop **licensing and managed settings**: use a managed and approved install.
- Cap WSL2 memory/CPU (`.wslconfig`) so an agent session plus image builds do not starve the machine.

**Worked example (both repos)**

```bash
task docker:check     # build the `test` stage: format/vet + tests inside Linux. Must pass.
task docker:up        # build runtime image, start the stack
curl http://localhost:8080/healthz      # {"status":"ok"}
curl -X POST localhost:8080/items/A/add -d '{"name":"Alpha","quantity":3}'
curl localhost:8080/items/A
task docker:down
```

Now watch the guardrails work. In a Claude Code session ask: "Start the service in Docker and add stock via curl." With the shipped `settings.json` the agent may run `task docker:up` only after your approval, and `curl` is denied outright (see the deny list), so a well-behaved agent will tell you it cannot call the API. That is the cue to design a deterministic `task docker:smoke` target for it instead of loosening permissions.

**Exercise:** (1) Write `task docker:smoke` that waits for health, adds stock, reads it back, and fails on mismatch. (2) Put the agent in a dev container (`templates/docker/devcontainer.json`), mount only the repo, and confirm it cannot read your home directory. (3) Break the build on purpose (introduce a CRLF `Dockerfile`) and see how `.gitattributes` prevents it.

**Done when:** `task docker:check` and `task docker:up` work on Docker Desktop, the agent can reach Docker only through allow-listed tasks, and CI builds the same `test` target.

---

### Step 5 · Explore before you change anything

**Goal:** use the agent as a read-only guide to unfamiliar code, and validate its understanding.

**Why it matters:** Most wasted agent effort comes from a wrong mental model. Exploration in plan/read-only mode is cheap and safe.

**What to do**

1. Use `templates/prompts/01-explore-codebase.md` in read-only mode.
2. Make it cite file paths for every claim. Spot-check three of them.
3. Ask for what is *surprising*: dead code, missing tests, inconsistent patterns.
4. Feed corrections back: if the agent misunderstood something, that misunderstanding belongs in `AGENTS.md` or `docs/architecture.md`.
5. For big repos, use read-only subagents or separate sessions for exploration so your main context stays clean.

**Worked example prompt (shelfstock):**

> Explore this repo read-only. Explain how a "remove stock" request flows from HTTP to the store, which errors can occur, and how each maps to an HTTP status. Cite files and functions.

Expected answer: `httpapi.remove` -> `decode` -> `Store.Remove` -> sentinel errors -> `writeErr` -> 400/404/409. If the agent says 500 for insufficient stock, your docs or the agent is wrong. Find out which.

**Exercise:** pick a real service you own. Run the explore prompt. List every inaccuracy in the agent's answer. Each one is a missing doc.

**Done when:** you trust the agent's summary of a module enough to plan changes from it.

---

### Step 6 · Spec, plan, implement

**Goal:** apply the core workflow to a real feature.

This is the central loop. Keep it identical across tools.

```
 ┌────────┐   ┌────────┐   ┌────────┐   ┌───────────┐   ┌────────┐   ┌────────┐
 │ SPEC   │ → │ PLAN   │ → │ SLICE  │ → │ VERIFY    │ → │ REVIEW │ → │ MERGE  │
 │ human  │   │ agent  │   │ agent  │   │ task check│   │ human  │   │ human  │
 │ + agent│   │ + human│   │ test-  │   │ + 2nd     │   │ + 2nd  │   │        │
 │ interv.│   │ approve│   │ first  │   │ agent     │   │ agent  │   │        │
 └────────┘   └────────┘   └────────┘   └───────────┘   └────────┘   └────────┘
   Step 6        Step 6       Step 7       Step 7         Step 8
```

**Why it matters:** Vague prompts produce plausible, wrong code at scale. A spec turns intent into testable criteria. A plan exposes misunderstandings before code exists. Correcting a plan is nearly free; correcting 800 lines is not.

#### 6.1 Spec

1. Skip the spec only for trivial changes (typo, one-line fix, obvious bug). Anything over ~1 hour deserves one.
2. Use `templates/prompts/02-write-spec.md`: the agent **interviews you** with clarifying questions, then drafts the spec in `docs/specs/<name>.md`.
3. Insist on **testable acceptance criteria**: Given/When/Then or concrete input/output tables. If you cannot write a test for it, it is not a requirement yet.
4. Always include **out of scope**. It is your scope-creep defence against an eager agent.
5. A human reads and approves the spec. This is the highest-value review in the whole process.

**Worked example:** `docs/specs/low-stock-alerts.md`. Notice:

- The behaviour table has concrete cases, including the zero-threshold "disabled" rule.
- The edge-case section pre-empts a real trap: `GET /items/low-stock` could be captured by `GET /items/{sku}`.
- "Out of scope" excludes notifications and persistence.

#### 6.2 Plan

1. Switch to plan mode / read-only. Use `templates/prompts/03-plan.md`.
2. Require: files to touch, ordered small steps (each buildable and testable), interfaces, risks, and what will *not* be done.
3. **Review the plan critically.** Look for: steps that are too big, missing tests, touching files outside scope, new dependencies, plan/spec conflicts.
4. Approve, or send it back with corrections. Do not let it start coding until you are satisfied.

**Worked example: a good plan for low-stock alerts**

| # | Step | Test that proves it |
|---|---|---|
| 1 | Add `Threshold` to `Item`, `ErrInvalidThreshold`, `Store.SetThreshold` | table test: ok / unknown SKU / negative |
| 2 | Add `Store.LowStock()` | test with A(3,thr 5), B(9,thr 5), C(thr 0) -> only A; plus empty returns `[]` not `nil` |
| 3 | HTTP: `POST /items/{sku}/threshold`, `GET /items/low-stock`, map new error to 400 | handler tests: 200/400/404 and route precedence |
| 4 | Update `docs/architecture.md` | `task check` still passes |

If an agent's plan merges steps 1-3 into one commit, push back: you want small reviewable slices.

#### 6.3 Implement

Continue with [Step 7](#step-7--test-first-slices-and-verification-loops).

**Exercise:** run 6.1 and 6.2 for the low-stock feature on your own, in your preferred tool. Then run 6.1 and 6.2 again in the *other* tool (Claude Code vs Codex) using the same prompts from `templates/prompts/`. Compare. Differences show where your written context is thin.

**Done when:** you have an approved spec and an approved plan with small, testable steps.

---

### Step 7 · Test-first slices and verification loops

**Goal:** implement the plan one slice at a time with the agent proving each slice.

**Why it matters:** Tests give the agent an objective success signal. Test-first prevents the classic failure mode where the agent writes code, then writes tests that merely confirm it.

**The loop for each slice** (`templates/prompts/04-implement-slice.md`):

1. Agent writes failing test(s) and shows they **fail for the right reason**.
2. Agent writes minimal code to pass.
3. Agent runs `task check` and fixes failures properly.
4. Agent reports: files changed, tests added, exact commands + results.
5. **You** skim the diff, then commit. One slice = one commit.
6. Clear or compact context before the next slice if the session is getting long.

**Rules to enforce**

- Never allow: skipping, deleting, or weakening tests or linters to get green.
- Bug fix = regression test that fails first (`templates/prompts/07-debug.md`).
- If the agent touches files outside the plan, stop and ask why.
- Run the thing. Compiling is not working. For services, hit an endpoint.

**Worked example: slice 1**

Prompt (condensed from the template):

> Implement step 1 of docs/specs/low-stock-alerts.md. Write failing tests first and show they fail. Then implement. Run `task check`. Stop after this step.

What good output looks like (Go): a table-driven `TestSetThreshold` with cases for success, unknown SKU, and negative threshold; a red run first (`undefined: SetThreshold`); then minimal implementation; then green `task check`. In .NET the same slice is an xUnit `[Fact]`/`[Theory]` set; the red run is a compile error (`'InventoryStore' does not contain a definition for 'SetThreshold'`), which counts as red.

**Reference gotchas hit while building the solution** (add such lines to `AGENTS.md` when they bite you):

| Gotcha | How it was caught | Rule |
|---|---|---|
| `LowStock()` returning a nil slice serialises to JSON `null`, not `[]` | A dedicated test `TestLowStockEmptyIsNotNil` | Return empty slices from list-style APIs |
| `GET /items/low-stock` vs `GET /items/{sku}` routing | A real HTTP smoke test showed the literal route wins | Go 1.22 mux prefers the more specific pattern; test routes, don't assume |
| Whitespace-sensitive scripted edit silently failed | Compiler error on the next run | Always run the build after an edit |
| **.NET:** `Item` is an immutable record, and `Add` builds a *new* `Item`, silently resetting `Threshold` to 0 | The spec's edge-case list plus test `Add_and_Remove_preserve_the_threshold` | Use `with` expressions or carry every field; write the "preserved across updates" criterion into the spec |
| **.NET:** `GET /items/low-stock` vs `/items/{sku}` | Test `Low_stock_route_is_not_captured_by_the_sku_route` | ASP.NET ranks literal segments above parameters, but prove it |
| **.NET:** malformed JSON returns 400 or throws depending on environment | `RouteHandlerOptions.ThrowOnBadRequest = false` in `Program.cs` plus status-code tests | Make behaviour explicit instead of depending on Development vs Production defaults |
| **.NET:** `dotnet format --verify-no-changes` fails on CRLF or re-spacing | `task lint`, `task docker:check` | `.gitattributes` + `.editorconfig`; run `task fmt` (the shipped hook does) |
| **.NET:** floating package version `10.0.*` in the test project | Reviewer checklist | Pin with `Directory.Packages.props` and lock files before production use |

**Exercise:** implement all four plan steps. Then compare to `examples/low-stock-reference-solution-go.patch`:

```bash
# Go track (use `git apply --check -p1 <file>` if `patch` is not installed, e.g. on Windows)
cd examples/shelfstock-go     && patch -p1 --dry-run < ../low-stock-reference-solution-go.patch
# .NET track (This is uncompiled reference; needs review: treat failures as part of the exercise)
cd examples/shelfstock-dotnet && patch -p1 --dry-run < ../low-stock-reference-solution-dotnet.patch
```

Both patches are *dry runs*: inspect them, do not apply them over your own work.

**Done when:** feature complete, each slice committed separately, `task check` green.

---

### Step 8 · Review and quality gates

**Goal:** add independent checks so one agent's blind spot is not the whole system's.

**Why it matters:** The author (human or AI) is the worst reviewer of their own work. AI makes it cheap to produce code and expensive to *read* it. Your bottleneck moves from writing to reviewing, so design for it.

**Gate stack**

| Gate | Who/what | Catches |
|---|---|---|
| 1. Local `task check` | Agent + author | Build, lint, test failures |
| 2. Independent AI review in a **fresh** session | A second agent with no authoring context (`templates/prompts/06-review.md`; Claude Code has `/review`-style commands; Codex has a review command) | Logic errors, missed edge cases, spec drift |
| 3. CI | Same `task check` + container build (`docker build --target test .`) + security scan + dependency/licence checks | Environment drift, vulnerable deps, leaked secrets |
| 4. Human review with the PR checklist | A person | Design, intent, business rules, "should we even do this?" |
| 5. Branch protection | Platform | Unreviewed merges |

**What to do**

1. Add `templates/github/pull_request_template.md` and `PR-CHECKLIST.md`. Include "how was this made" (human / AI-assisted / mostly AI) so reviewers adjust attention.
2. Add `templates/github/ci.yml` (language-neutral skeleton) running `task check` plus your approved secret scanner.
3. Use a **different** session (ideally a different model or tool) to review. Authoring and reviewing in one context shares the same blind spots.
4. Review for the failure patterns specific to AI code:
   - Hallucinated APIs, flags, or package functions.
   - Tests that cannot fail (tautologies, over-mocking).
   - Unrequested changes: reformatting, config edits, deleted tests.
   - Plausible-but-wrong logic: off-by-one, swapped arguments, wrong units/timezones.
   - Duplicated helpers instead of reuse.
5. **Cap PR size** (~400 lines). Large agent diffs are where review quality collapses. Ask the agent to split or stack PRs.
6. Mark AI-assisted commits (optional but recommended): a trailer such as `Assisted-by: Claude Code` or `Assisted-by: Codex`.

**Worked example:** after slice 3, open a *fresh* session and run the review prompt on your branch. A typical valid finding: "`threshold` of `0` conflates 'disabled' with 'zero', which the spec allows, but the JSON field name `threshold` is not documented in the API docs." Then run the human PR checklist yourself.

**Exercise:** deliberately introduce a bug (e.g. change `<=` to `<` in `LowStock`) on a branch and see whether (a) tests, (b) AI review, or (c) you catch it. Learn which gate is weakest.

**Done when:** no change can reach `main` without passing CI, an independent review, and a named human approver.

---

### Step 9 · Build reusable assets

**Goal:** turn good prompts and workflows into shared team assets.

**Why it matters:** Consistency is what makes results predictable across engineers. The asset library is where team's accumulated practice lives.

**The asset ladder, from most portable to most tool-specific**

| Asset | What it is | Portability |
|---|---|---|
| Prompt files | Markdown in `docs/prompts/` or `templates/prompts/` that you paste or reference | Any tool |
| `AGENTS.md` rules | Always-on guidance | Any tool |
| Scripts / Make targets | Deterministic automation agents can call | Any tool |
| Custom slash commands | Saved prompts invoked as `/name` (Claude Code: `.claude/commands/*.md`) | Per tool |
| Skills | Folders of instructions + scripts that load on demand (Claude Code skills; Codex supports skills too, confirm current layout) | Per tool, converging format |
| Subagents | Specialised agents with their own context and tool limits (e.g. read-only reviewer) | Per tool |
| MCP servers | Standard protocol for exposing tools/data (issue tracker, docs, DB read-replica) to agents | Works across Claude, Codex, Copilot |
| Hooks | Deterministic scripts on lifecycle events | Per tool |

**Guidance**

1. **Start with prompt files and Make targets.** They work everywhere and teach you what to automate. Promote to commands/skills once a prompt has been used ~5 times.
2. Prefer **deterministic code over instructions** when possible: a script that checks layering beats a paragraph asking the agent to respect layering.
3. **MCP servers need a security review** (Step 11): they expand what the agent can read and do. Prefer read-only. Pin versions. Keep an approved list.
4. Keep assets in the repo and review them like code. Each should have a one-line "when to use it".

**Worked example:** `templates/claude/commands/spec.md` wraps the spec prompt as `/spec <description>`. The same prompt text lives in `templates/prompts/02-write-spec.md` so a Codex user pastes it, or you wrap it in whatever custom-prompt/skill mechanism your Codex version supports. **One prompt, two wrappers.**

**Exercise:** take your most frequent prompt this week, turn it into a prompt file with a "when to use" header, and wire it as a command in Claude Code. Have a teammate use it in Codex via copy-paste. Ask them what was unclear.

**Done when:** your top 5 recurring tasks each have a shared prompt file.

---

### Step 10 · Manage context and long-running work

**Goal:** keep agent quality high over long tasks and across sessions and tools.

**Why it matters:** Context windows fill; quality degrades as they do. Long sessions accumulate stale assumptions and dead ends.

**Techniques**

| Technique | When |
|---|---|
| **One task per session.** Clear/restart between unrelated tasks | Always |
| **Compact** the conversation at natural breakpoints (`/compact` in Claude Code) | Mid-task when context is long but relevant |
| **Handoff note** (`templates/prompts/10-session-handoff.md`) | End of day, switching tools, or when the session feels "off" |
| **Plan in a file**, not just in chat (`docs/specs/` + checklists) | Anything multi-step. The file survives a context reset |
| **Subagents / separate sessions** for exploration and review | Keep the main thread clean |
| **Git worktrees** for parallel agents on separate branches | Running 2-3 independent tasks at once without collisions |
| **Frequent commits** | Easy rollback of a bad agent turn |

```bash
# Two independent agent sessions on two branches, no checkout collisions
git worktree add ../shelfstock-lowstock -b feat/low-stock
git worktree add ../shelfstock-healthz  -b feat/healthz
# start one agent in each directory
```

**Signs a session has gone bad:** it repeats itself, forgets rules from `AGENTS.md`, edits unrelated files, or "fixes" the same bug three ways. Stop. Write a handoff note. Start fresh. This is faster than arguing.

**Worked example:** halfway through the low-stock plan, run the handoff prompt, open a brand-new session (try the *other* tool), and tell it to read the note and continue with the next step. A good scaffold makes this nearly seamless. That is the test of whether your files, not your chat history, hold the project's state.

**Exercise:** do the Claude -> Codex handoff on shelfstock. List what the new agent got wrong. Fix the root cause in `AGENTS.md` or the handoff template.

**Done when:** you can switch tools mid-feature with no loss of continuity.

---

### Step 11 · Security, IP and compliance

**Goal:** a policy and technical controls that satisfy Security, Legal and Privacy, and let engineers move quickly inside them.

**The risk map**

| Risk | Example | Control |
|---|---|---|
| **Data leakage** | Customer PII or secrets sent to a model | Data classification; deny-lists; synthetic test data; enterprise plans with no-training and defined retention |
| **Prompt injection** | A web page, issue, or dependency README tells the agent to exfiltrate secrets or run commands | Treat tool/web content as data; restrict network and credentials; require approval for side-effects; least privilege |
| **Excessive agency** | Agent runs destructive commands or pushes to main | Sandboxes, deny rules, branch protection, no prod credentials |
| **Supply chain** | Agent adds an unvetted or typosquatted dependency | "Ask first" rule for dependencies; licence + CVE scanning in CI; lockfile review |
| **Insecure code** | Plausible but vulnerable code (injection, authz gaps) | Security review prompt; SAST in CI; mandatory expert review in sensitive areas |
| **IP / licence contamination** | Model reproduces licensed code | Vendor indemnity terms; duplicate-detection filters; licence scanning; policy against pasting third-party proprietary code |
| **Compliance gaps** | No audit trail; regulated data in prompts | Tool inventory; logging; DPAs; mapped controls (PIPEDA/GDPR/PCI/SOX as applicable) |
| **Docker access / container escape** | Agent with the Docker socket runs a privileged container or mounts the host filesystem | No socket for agents; allow-list `task docker:*` only; disposable dev containers; deny raw `docker run`; see Step 4b |
| **MCP / plugin risk** | Malicious or over-privileged tool server | Approved list; read-only by default; pinned versions; review |

**What to do**

1. Start from `templates/docs/ai-security-and-compliance-policy.md`. Have Security, Legal and Privacy co-author; this handbook is not legal advice.
2. For each approved tool, record in a table: plan/tier, **training on your data (must be off)**, retention, data residency, certifications, output-ownership/indemnity terms, admin/audit features, owner.
3. Implement controls as **code**, not just policy: deny rules, sandbox config, CI scanners, branch protection, CODEOWNERS for sensitive paths.
4. Decide your "no-go" data classes (typically: customer PII, payment data, credentials, production data) and enforce them technically.
5. Define the incident procedure: revoke credentials, notify security, record the new rule.
6. Retail/e-commerce note: if payment data is in scope, check PCI-DSS implications of any code path that touches cardholder data *before* allowing agent access to it. Typically keep that code in a restricted repo/area with extra review.

**Worked example:** the shelfstock `AGENTS.md` "Security & data handling" section, `.claude/settings.json` deny rules, and CI permissions (`contents: read`) show the three layers (instruction, enforcement, platform) for a trivial service.

**Exercise:** threat-model the shelfstock HTTP API with `templates/prompts/09-security-review.md`. Which findings are real (no auth, no rate limit, unbounded `GET /items`)? Add the real ones to the architecture doc's "Known debt".

**Done when:** every approved tool has a completed vendor-risk row, controls are enforced in code, and the policy has Security/Legal sign-off.

---

### Step 12 · Measure what matters

**Goal:** know whether AI assistance is helping, and catch harm early.

**Why it matters:** Perceived productivity is unreliable. Leaders who only count "lines generated" or "suggestions accepted" will optimise the wrong thing. Measure outcomes your business cares about, and watch quality as hard as speed.

**Suggested scorecard** (measure before rollout for a baseline)

| Category | Metric | Notes |
|---|---|---|
| **Speed** | PR cycle time (open -> merge); lead time for changes | DORA-style |
| **Throughput** | Deployment frequency; PRs merged per engineer | Use trends, not rankings |
| **Quality** | Change failure rate; escaped defects; reverts/hotfixes; flaky test rate | Must not worsen |
| **Review health** | Review time; PR size; rework requested per PR | AI often inflates PR size: watch it |
| **Security** | Findings per PR; vulnerable dependency count; secrets caught pre-merge | |
| **Adoption** | Weekly active users; sessions per user; share of PRs flagged AI-assisted | Adoption is an input, not the goal |
| **Developer experience** | Short quarterly survey: trust, friction, time saved, skills concerns | Qualitative signal matters |
| **Cost** | Licence + usage spend per team; value per spend | |

**Rules for using metrics**

1. Never use them for individual performance ranking. People will game them or stop being honest.
2. Compare like-for-like (same team, similar work) before/after, and treat early numbers with scepticism.
3. A speed gain with a rising change-failure rate is a loss.
4. Review the scorecard monthly; adjust practices (e.g. tighter PR size caps) in response.

**Worked example:** for your proposal, define "reduce cycle time" precisely: median PR open-to-merge, per team, 8 weeks pre-pilot vs 8 weeks during pilot, with change-failure rate and PR size alongside.

**Exercise:** write a one-page measurement plan for your pilot team with a baseline period, five metrics, and the threshold at which you would pause the pilot.

**Done when:** you have a baseline captured *before* the pilot begins.

---

### Step 13 · Roll out to a team

**Goal:** a phased, low-risk adoption plan for a team new to coding agents.

**Why it matters:** The failure modes of rollout are cultural and process-related far more than technical: unequal skill, no standards, fear, and over-trust.

This phase order matches the structure of a typical proposal for Codex adoption, ordered by increasing autonomy and risk:

| Phase | Name | What agents do | Human role | Exit criteria |
|---|---|---|---|---|
| **0** | Foundations | Nothing yet | Fix `task check`, write `AGENTS.md`, set up CI/gates, approve tools and policy | 2-3 pilot repos are agent-ready |
| **1** | AI review in the pipeline | Review PRs and builds (advisory comments only) | Humans still review and merge; triage AI comments | Useful-comment rate acceptable; no process friction |
| **2** | Agent-written tests | Write unit tests, characterization tests for under-tested code | Review tests for meaning, not just coverage | Coverage up, no tautological tests, reviewers confident |
| **3** | Engineer-directed agents with a harness | Implement small features/fixes from specs; verify with `task check` | Engineers write specs, approve plans, review diffs | Cycle time improves with stable quality metrics |
| **4** | Workflow automation | Triggered by tickets/CI; produce PRs unattended for narrow, low-risk classes (dependency bumps, doc fixes, flaky test triage) | Humans review every PR; scope allow-listed | Low revert rate; clear kill switch |

**Operating the rollout**

1. **Pick champions.** 1-2 enthusiastic, respected engineers per team. They test the scaffold first and teach peers.
2. **Pilot with 2-3 repos and ~5-8 engineers** for 6-8 weeks. Include one legacy repo, since that is where you learn most.
3. **Train everyone with the same short curriculum** (use this handbook's Steps 5-8 as workshop modules). Hands-on, using shelfstock, then their own code.
4. **Standardise one setup per tool**: the same `AGENTS.md` skeleton, the same permission baseline, the same prompt library. Make the golden path easy.
5. **Run a weekly 30-minute "what worked, what burned us" session.** Convert every lesson into an `AGENTS.md` line, prompt update, or checklist item.
6. **Set explicit expectations about accountability**: the human owns the code, can explain every line, and discloses AI assistance in the PR.
7. **Address the human side openly**: junior engineers must still learn fundamentals (consider "explain this diff" exercises and pairing); seniors shift toward specs, review and architecture; nobody is being measured on AI usage volume.
8. **Have a pause/rollback plan**: define the metrics thresholds and incidents that trigger a pause.
9. **Share a "known limits" list** with the team and refresh it monthly.

**Windows and Docker Desktop readiness (check before Phase 1)**

1. Standard install for every laptop: Git, Task, SDK(s), Docker Desktop (licensed, managed), WSL2 if used. Script it (winget/choco/Intune).
2. Decide the agent environment per tool (native Windows vs WSL) and document it once; keep repo and agent in the same environment.
3. `.gitattributes`, `.editorconfig`, and the `docker:check` job in every pilot repo before the first agent session.
4. Choose operating model A (host-native + Docker parity) for the pilot; move to model B (dev containers) before enabling any higher-autonomy mode.
5. Add `Taskfile.yaml`, `Dockerfile`, `compose.yaml` and CI to CODEOWNERS.

**Common rollout traps**

| Trap | Counter |
|---|---|
| Handing out licences with no scaffold | Do Phase 0 first |
| Mandating usage metrics | Measure outcomes, not usage |
| Huge AI-generated PRs overwhelming reviewers | PR size caps; stacked PRs |
| Two teams, two incompatible setups | One shared template repo (this kit) |
| Security finds out afterwards | Involve Security in Phase 0 |
| "Works on my laptop" (Windows vs macOS vs CI) | One `docker:check` target and one Taskfile vocabulary; `.gitattributes` |
| Skill atrophy / junior concerns | Intentional practice: write-by-hand sessions, explain-the-diff reviews |
| Overclaiming early wins | Baseline first; report medians and quality alongside speed |

**Worked example:** convert this kit into a template repository. The team clones it, replaces the shelfstock example with their repo's content, and fills in the policy. Then run one live 90-minute workshop: all engineers apply Steps 5-7 to the low-stock feature in pairs, one pair using Claude Code, one using Codex.

**Exercise:** draft your own Phase 0 checklist and the pilot's success/pause criteria, then share them with Security and a peer director for challenge.

**Done when:** you have a written rollout plan with owners, dates, metrics, and pause criteria, and Security and Legal have signed off on Phase 0.

---

## 5. Cross-tool mapping

The same concepts under different names. **Verify details against current docs before standardising.**

| Concept | Claude Code | Codex | GitHub Copilot |
|---|---|---|---|
| **Project instructions** | `CLAUDE.md` (import `@AGENTS.md`) | `AGENTS.md` (native; nested overrides; size cap `project_doc_max_bytes`, default 32 KiB) | `.github/copilot-instructions.md`; agent also supports `AGENTS.md` |
| **User-level instructions** | `~/.claude/CLAUDE.md` | `~/.codex/AGENTS.md` | Personal / org-level custom instructions |
| **Permissions** | `.claude/settings.json` `permissions.allow/deny`; modes (plan, accept edits, bypass) | `approval_policy`, `sandbox_mode` (`read-only` / `workspace-write` / `danger-full-access`) in `~/.codex/config.toml` or project `.codex/config.toml` | Platform and org policies; agent runs in controlled environment |
| **Network control** | Deny rules for fetch tools; sandboxing | `[sandbox_workspace_write] network_access` (default off) | Firewall settings for the cloud agent |
| **Reusable prompts** | `.claude/commands/*.md` slash commands | Custom prompts / skills (confirm current mechanism) | Prompt files (`.github/prompts/`) |
| **Skills** | `.claude/skills/` | Skills supported (confirm layout) | Agent skills (emerging) |
| **Subagents** | `.claude/agents/*.md` | Via task delegation/new sessions | Custom agents |
| **Hooks** | `settings.json` `hooks` (PreToolUse, PostToolUse, etc.) | Notify hooks / rules (confirm) | Limited |
| **MCP** | `.mcp.json` / `claude mcp add` | `[mcp_servers.*]` in `config.toml` | MCP in agent mode / org config |
| **Non-interactive / CI** | `claude -p "..."` (headless) | `codex exec "..."` | Coding agent via issues/PRs |
| **Code review** | Review commands / GitHub integration | `/review` and GitHub review integration | Copilot code review on PRs |
| **Plan / read-only** | Plan mode | Read-only sandbox | Plan-first behaviour in agent mode |
| **Context management** | `/compact`, `/clear`, `/resume` | New session / resume, compact | Chat resets |
| **Windows** | Native Windows and WSL are both supported (confirm current requirements) | WSL commonly recommended; native Windows support has been experimental (confirm) | Works in VS Code / Visual Studio natively; cloud agent is Linux |
| **Container isolation** | Run inside a dev container/VM; permissions rules for `docker` commands | Sandbox modes plus container/VM isolation | Cloud agent runs in an isolated environment |
| **Init draft** | `/init` | `/init` | Generate instructions from repo |

**What this table tells you:** only the right-hand columns of *permissions, hooks and commands* differ. The left-hand pillars (AGENTS.md, `task check`, spec/plan/slice process, CI, policy) are identical. **Keep your investment in the portable layer.**

**Sharing one prompt across tools:** keep the canonical prompt text in `templates/prompts/*.md`. In Claude Code, wrap as a slash command. In Codex, either paste, reference the file ("follow `docs/prompts/06-review.md`"), or wrap in the supported custom-prompt/skill mechanism. Never fork the prompt text per tool.

---

## 6. Porting to other languages

The method does not change. Only the contents of `check` and the conventions section do. Fill this table in for each repo you onboard:

| Language | `build` | `test` | `lint` / `fmt` | Typical `AGENTS.md` gotchas |
|---|---|---|---|---|
| **Go** (`examples/shelfstock-go`) | `task build` -> `go build ./...` | `task test` -> `go test -race ./...` (race needs cgo; not on plain Windows) | `go vet`, `gofmt -l`, `golangci-lint` | Error wrapping with `%w`; table-driven tests; nil vs empty slice in JSON; layering by package convention |
| **TypeScript / Node** | `npm run build` | `npm test` (Vitest/Jest) | `eslint`, `prettier`, `tsc --noEmit` | ESM vs CJS; strict mode; do not add `any`; lockfile discipline |
| **Python** | `python -m build` | `pytest` | `ruff`, `black`/`ruff format`, `mypy` | Virtualenv/uv usage; type hints; avoid mutable defaults |
| **Java / Kotlin** | `./gradlew build` / `mvn -q package` | `./gradlew test` | `spotless`, `checkstyle`, `detekt` | Spring wiring; null-safety; JPA lazy loading |
| **C# / .NET** (`examples/shelfstock-dotnet`) | `task build` -> `dotnet build` (warnings as errors) | `task test` -> `dotnet test` (xUnit, `WebApplicationFactory`) | `dotnet format --verify-no-changes` + `.editorconfig` + analyzers | Async all the way; nullable reference types; records are immutable (`with`); LF endings; floating package versions; layering enforced by project references |
| **Rust** | `cargo build` | `cargo test` | `cargo clippy -- -D warnings`, `cargo fmt --check` | Ownership-driven refactors can sprawl; keep slices small |
| **Ruby** | `bundle exec rake` | `bundle exec rspec` | `rubocop` | Metaprogramming hides call sites; rely on tests |
| **Legacy / no tests** | whatever builds | **write characterization tests first** (prompt 05) | add a formatter in its own PR | Document landmines heavily; require human review on every change |

**Onboarding a codebase with no tests:** the first agent task is not a feature. It is "add characterization tests around the module we are about to change" (`templates/prompts/05-write-tests.md`), then refactor. Without tests the agent has no steering wheel.

#### 6.1 The two example repos side by side

| Aspect | Go (`shelfstock-go`) | .NET (`shelfstock-dotnet`) |
|---|---|---|
| Layout | `cmd/`, `internal/inventory`, `internal/httpapi` | `src/ShelfStock.Core`, `src/ShelfStock.Api`, `tests/ShelfStock.Tests`, `ShelfStock.slnx` |
| Layering rule enforced by | Convention in `AGENTS.md` + review | **Compiler** (project references): Core cannot reference Api |
| Errors | Sentinel errors + `errors.Is` | Exception hierarchy (`InventoryException`) mapped in one place |
| "Warnings are errors" | `go vet` in `lint` | `TreatWarningsAsErrors` + nullable + analyzers in `Directory.Build.props` |
| Format gate | `gofmt -l` | `dotnet format --verify-no-changes` |
| Tests | table-driven, `httptest` | xUnit `[Theory]`, `WebApplicationFactory<Program>` |
| Toolchain pin | `go.mod` | `global.json` + `Directory.Build.props` |
| Docker `test` stage | `go vet` + `go test -race` | `dotnet format` check + `dotnet test` |
| Runtime image | static binary on `alpine`, non-root | `aspnet` image, non-root (`$APP_UID`), `tmpfs /tmp` |
| Windows pitfall | `-race` needs a C compiler | CRLF breaks `dotnet format` |

**What stays identical:** `AGENTS.md` structure, the `task` vocabulary, `docs/` layout, spec, prompts, PR checklist, CI shape. **That is the point of the kit**: only the contents of a few files change per language.

**Polyglot monorepos:** root `AGENTS.md` with the shared vocabulary (`task check` at the root fans out), plus nested `AGENTS.md` per service with language-specific commands. With Task, a root `Taskfile.yaml` can pull in per-service files: `includes: { api: { taskfile: ./services/api, dir: ./services/api } }` and a root `check` that calls `api:check` and `web:check`.

---

## 7. Anti-patterns and troubleshooting

| Anti-pattern | Symptom | Fix |
|---|---|---|
| **Vibe coding in production** | Large unreviewed diffs, nobody understands the code | Spec/plan/slice; PR size caps; "explain every line" accountability |
| **Giant prompt, no spec** | Plausible code that does the wrong thing | Write the spec; interview first |
| **Prompting instead of documenting** | You repeat the same corrections every session | Put the correction in `AGENTS.md` |
| **Bloated `AGENTS.md`** | Rules ignored or truncated | Prune; link to docs; split into nested files |
| **No fast check** | Agent says "done" without proof | Step 1 |
| **Agent grades its own work** | Same blind spots in author and reviewer | Fresh-session or different-tool review |
| **Weakened tests to go green** | Skipped tests, loosened assertions | Hard "Never" rule + reviewer checklist item + CI guard against skipped tests |
| **Full-access mode on a laptop** | Unintended file/system damage | Sandboxed or disposable environments only |
| **Marathon sessions** | Drift and repeated mistakes | Handoff notes; fresh sessions |
| **One mega-PR** | Rubber-stamp review | Stacked small PRs |
| **Trusting confident prose** | Hallucinated APIs | Compile, run, read docs for new APIs |
| **Agent given the Docker socket / raw `docker run`** | Host-level compromise from one bad command or injected instruction | Allow-list `task docker:*` only; dev containers; Step 4b |
| **Editing Taskfile/Dockerfile/CI without review** | Agent widens its own permissions by wrapping commands in tasks | CODEOWNERS + "Ask first" rule |
| **CRLF line endings from Windows editors** | Scripts and `dotnet format` fail only in containers/CI | `.gitattributes` `eol=lf`, `.editorconfig` |
| **Metrics theatre** | "AI wrote 40% of our code" as the success metric | Outcome and quality metrics |

**Troubleshooting a wrong or low-quality result**

1. Was the task small and specific? If not, split it.
2. Was the relevant context available (files, docs, spec)? If not, add it to the repo, not the chat.
3. Could the agent verify its own work? If not, add or fix the check.
4. Is the session long or confused? Handoff and restart.
5. Is the model wrong for the task? Use a stronger model for planning and review, a faster one for mechanical edits.
6. Still bad? Do that part by hand. That is allowed and sometimes optimal.

---

## Appendices

### Appendix A · Quick-start checklist (per repository)

- [ ] `task check` (or equivalent) exists, is fast, and is green
- [ ] Task version pinned in CI/devcontainer; `task --list` smoke-tested
- [ ] `Dockerfile` with `test` + `runtime` stages; `task docker:check` green; `compose.yaml` runs locally on Docker Desktop
- [ ] `.gitattributes` (`eol=lf`) and `.editorconfig` committed
- [ ] `Taskfile.yaml`, `Dockerfile`, `compose.yaml`, CI in CODEOWNERS
- [ ] CI runs the same command
- [ ] `AGENTS.md` committed (< 150 lines), `CLAUDE.md` imports it
- [ ] `docs/architecture.md` and `docs/adr/` exist
- [ ] `.claude/settings.json` deny rules; Codex sandbox/approval set
- [ ] `.env`/secrets denied to agents; secret scanner in CI
- [ ] PR template + `PR-CHECKLIST.md`
- [ ] Branch protection: required review and required checks
- [ ] Prompt library in the repo (`docs/prompts/` or copied from `templates/prompts/`)
- [ ] Spec template in `docs/specs/_template.md`
- [ ] Policy acknowledged by team

### Appendix B · The daily loop (one page)

```
1. Start fresh session. Pick ONE task.
2. Is it > 1 hour?  yes -> Spec (interview) -> Human approves
3. Plan mode -> Plan -> Human approves
4. For each slice: failing test -> minimal code -> task check -> commit
5. Fresh-session AI review -> fix
6. Human self-review with PR checklist -> open PR (disclose AI assistance)
7. Human reviewer approves -> merge
8. Anything surprising?  -> add a line to AGENTS.md
```

### Appendix C · Further reading

Verify these are current before relying on specifics.

- AGENTS.md open format: https://agents.md
- Claude Code docs (memory/CLAUDE.md, settings, hooks, slash commands, subagents, skills, MCP): https://docs.claude.com/en/docs/claude-code/overview
- Codex docs (AGENTS.md, config, sandbox/approvals, skills, MCP): https://developers.openai.com/codex
- Codex sample configuration (keys referenced in `templates/codex/config.toml`): https://developers.openai.com/codex/config-sample
- GitHub Copilot custom instructions: https://docs.github.com/en/copilot/how-tos/configure-custom-instructions
- Architecture Decision Records: https://adr.github.io
- C4 model: https://c4model.com
- DORA metrics: https://dora.dev
- OWASP Top 10 for LLM Applications (prompt injection, excessive agency, supply chain): https://genai.owasp.org
- Model Context Protocol: https://modelcontextprotocol.io
- Task (Taskfile) docs: https://taskfile.dev
- Docker Desktop, Compose, multi-stage builds, Compose Watch: https://docs.docker.com
- .NET SDK, `dotnet format`, `global.json`, minimal APIs, integration tests: https://learn.microsoft.com/dotnet
- Dev Containers specification: https://containers.dev

### Appendix D · Glossary

| Term | Meaning |
|---|---|
| **Agent** | Model + tools + loop that can act on a codebase |
| **Context window** | The text the model can consider at once; finite |
| **AGENTS.md** | Open, tool-neutral instruction file for coding agents |
| **Plan mode** | Read-only mode where the agent proposes but cannot change files |
| **Sandbox** | Restricted environment limiting filesystem/network access |
| **Hook** | Deterministic script triggered by an agent lifecycle event |
| **MCP** | Model Context Protocol; standard way to connect tools/data to agents |
| **Skill** | Packaged instructions/scripts an agent loads on demand |
| **Subagent** | A helper agent with its own context and constraints |
| **Prompt injection** | Malicious instructions hidden in content the agent reads |
| **ADR** | Architecture Decision Record |
| **Characterization test** | Test that records current behaviour of legacy code |
| **Taskfile** | `Taskfile.yaml` read by the Task runner; defines the shared command vocabulary |
| **Docker Desktop** | Local Docker engine + UI for Windows/macOS; commonly uses WSL2 on Windows |
| **Compose** | `compose.yaml`-defined multi-container local stack (`docker compose up`) |
| **Dev container** | A container used as the development environment (toolchain + sandbox) |
| **Spec-driven development** | Writing testable requirements before code, and directing agents from them |

### Appendix E · Template index

| File | Use |
|---|---|
| `templates/AGENTS.md` | Starter instruction file |
| `templates/Taskfile.yaml` | Language-neutral Taskfile skeleton (with `docker:*` tasks) |
| `templates/docker/` | `Dockerfile.go`, `Dockerfile.dotnet`, `compose.yaml`, `devcontainer.json`, `dockerignore.template` |
| `templates/dotnet/` | `.editorconfig`, `Directory.Build.props`, `global.json` starters for .NET |
| `templates/.gitattributes.template` | LF line endings (rename to `.gitattributes`) |
| `templates/CLAUDE.md` | Claude Code shim importing `AGENTS.md` |
| `templates/claude/settings.json` | Claude Code permissions + format hook |
| `templates/claude/commands/spec.md`, `review.md` | Example slash commands |
| `templates/codex/config.toml` | Codex sandbox and approval baseline |
| `templates/github/copilot-instructions.md` | Copilot pointer to `AGENTS.md` |
| `templates/github/pull_request_template.md` | PR template |
| `templates/github/ci.yml` | Language-neutral CI skeleton |
| `templates/PR-CHECKLIST.md` | Full review checklist |
| `templates/prompts/01..10-*.md` | Prompt library (explore, spec, plan, implement, tests, review, debug, refactor, security, handoff) |
| `templates/docs/architecture.md` | Architecture doc skeleton |
| `templates/docs/adr-template.md` | ADR template |
| `templates/docs/spec-template.md` | Feature spec template |
| `templates/docs/ai-security-and-compliance-policy.md` | Starter policy for Security/Legal review |
