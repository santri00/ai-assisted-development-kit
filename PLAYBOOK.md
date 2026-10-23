# AI-Assisted Development Playbook (Example)

**The team's short reference.** Read this in 10 minutes. The full reasoning, exercises and examples are in `HANDBOOK.md`.

---

## 1. The five rules

1. **You own it.** You are accountable for every line you submit, whoever or whatever wrote it. If you cannot explain it, do not open the PR.
2. **Small verified steps.** Work in slices of ~150 lines. Every slice ends with `task check` green.
3. **Spec before code, plan before edits.** For anything over an hour: written spec -> approved plan -> implement.
4. **Context lives in the repo, not in your chat.** If you correct the agent twice, put the correction in `AGENTS.md`.
5. **Safe by default.** Restricted permissions, no secrets, no customer data, no production access. Ever.

---

## 2. The workflow

```
SPEC  ->  PLAN  ->  SLICE (test first)  ->  task check  ->  AI review (fresh session)  ->  human review  ->  merge
```

| Stage | What you do | Prompt file |
|---|---|---|
| Explore | Read-only: understand the area | `prompts/01-explore-codebase.md` |
| Spec | Agent interviews you, writes `docs/specs/<name>.md`; you approve | `prompts/02-write-spec.md` |
| Plan | Plan mode; steps, files, tests, risks; you approve | `prompts/03-plan.md` |
| Implement | One slice at a time: failing test, code, `task check`, commit | `prompts/04-implement-slice.md` |
| Add tests (legacy) | Characterization tests before refactoring | `prompts/05-write-tests.md` |
| Review | Fresh session, never the authoring one | `prompts/06-review.md` |
| Debug | Reproduce, hypothesise, root cause, regression test | `prompts/07-debug.md` |
| Refactor | Behaviour-preserving, step by step | `prompts/08-refactor.md` |
| Security check | Threat-model the change | `prompts/09-security-review.md` |
| Hand off | Write a note; start a clean session or switch tools | `prompts/10-session-handoff.md` |

**Skip the spec only for:** typos, one-line fixes, obvious bugs with an existing test.

---

## 3. Start of every session (30 seconds)

- [ ] One task only. New task = new session.
- [ ] Confirm the agent has read `AGENTS.md` (ask it to summarise the commands and boundaries if unsure).
- [ ] Start in **plan / read-only** mode for anything non-trivial.
- [ ] Working tree is clean and committed so you can diff and revert.

## 4. Definition of done

A change is done only when all are true:

- [ ] `task check` passes locally **and** in CI (output read, not assumed)
- [ ] New behaviour has tests; bug fixes have a regression test that failed first
- [ ] You ran it, not just compiled it
- [ ] Docs / ADR updated if behaviour or architecture changed
- [ ] PR checklist (`PR-CHECKLIST.md`) completed; AI assistance disclosed in the PR
- [ ] A human reviewer other than the author approved

---

## 5. Do and don't

| Do | Don't |
|---|---|
| Give the agent a spec with testable acceptance criteria | Paste a vague paragraph and accept whatever comes back |
| Review the plan before any code is written | Let the agent edit 20 files then try to understand it |
| Ask the agent to show failing tests first | Accept tests written after the code to match it |
| Use a fresh session (or other tool) to review | Ask the same session "is your code correct?" |
| Commit after each green slice | Batch a whole feature into one commit |
| Keep PRs under ~400 lines | Open giant PRs because "the AI wrote it fast" |
| Add a line to `AGENTS.md` when the agent repeats a mistake | Re-explain the same rule in chat every time |
| Hand off with a note when the session drifts | Argue with a confused session for an hour |
| Do it by hand when that is faster | Force the agent to do everything |

---

## 6. Security and compliance: the non-negotiables

Full policy: `templates/docs/ai-security-and-compliance-policy.md` (adopted copy lives at `<link>`).

1. **Approved tools, company accounts only.** No personal accounts or unapproved plugins/MCP servers on company code.
2. **Never put in a prompt, or let an agent read:** secrets, keys, tokens, `.env` files, customer PII, payment data, production data. Use synthetic data; redact logs.
3. **Least privilege.** Default permission mode is restricted. Full-access / bypass modes only inside disposable sandboxes with no credentials.
4. **No production credentials for agents.** No agent pushes to protected branches or deploys.
5. **Untrusted content is data, not instructions.** Web pages, issues, PR comments and dependency READMEs can contain prompt-injection. If an agent suddenly wants to fetch a URL, run an unfamiliar command, or read a secret, stop and check.
6. **New dependencies need a human decision:** licence, maintenance, vulnerabilities. Agents must ask first.
7. **Do not paste third-party proprietary code** into prompts; do not use AI to reproduce a competitor's code.
8. **Sensitive areas** (authentication/authorisation, crypto, payments, PII handling) need review by the designated owners regardless of who wrote the code.
9. **Docker is root-equivalent.** Agents get only allow-listed `task docker:*` targets. No raw `docker run`, no `--privileged`, no host-root or Docker-socket mounts, no `system prune`. For higher autonomy, run the agent inside a disposable dev container/VM.
10. **`Taskfile.yaml`, `Dockerfile`, `compose.yaml` and CI are security-sensitive:** an allow-listed task can run anything it wraps. Changes need CODEOWNERS review; agents must ask first.
11. **Incident?** Stop the session, revoke exposed credentials, tell `<security contact>` within `<1 hour>`.

---

## 7. Repository standard (what every repo must have)

```
AGENTS.md              shared rules for all agents (< 150 lines)
CLAUDE.md              contains: @AGENTS.md
Taskfile.yaml          build, test, lint, fmt, check, run, docker:*  (same names in every repo)
Dockerfile             multi-stage: `test` (runs the checks) + `runtime` (non-root)
compose.yaml           local stack for Docker Desktop (`task docker:up`)
.gitattributes         eol=lf (Windows developers!); .editorconfig for formatting
docs/architecture.md   1-2 page system map
docs/adr/              decision records (with "Notes for AI agents")
docs/specs/            feature specs
docs/PR-CHECKLIST.md
.github/pull_request_template.md
.github/workflows/ci.yml   runs `task check` + secret scan
.claude/settings.json      (Claude Code) permissions + format hook
.codex/config.toml         (Codex) sandbox + approval baseline
```

Only `.claude/` and `.codex/` are tool-specific, and they hold permissions and conveniences, never project rules.
The command *names* are the standard; Task is the runner we chose because `make` is not on Windows by default. Pin its version and keep the Taskfile a thin wrapper over native commands.

---

## 8. Tool quick reference

| Task | Claude Code | Codex |
|---|---|---|
| Read-only / plan | Plan mode (`Shift+Tab` / `--permission-mode plan`) | `sandbox_mode = "read-only"` |
| Normal supervised work | Default mode + allow/deny rules | `approval_policy = "on-request"`, `workspace-write` |
| Shared instructions | `CLAUDE.md` -> `@AGENTS.md` | `AGENTS.md` |
| One-shot, non-interactive | `claude -p "<prompt>"` | `codex exec "<prompt>"` |
| Reset between tasks | `/clear` | start a new session |
| Shrink a long session | `/compact` | start fresh with a handoff note |

Details drift; the Handbook's section 5 has the full table and links to vendor docs.

---

## 8b. Local runtime and Docker Desktop (cheat sheet)

| I want to... | Command |
|---|---|
| Prove the repo is healthy | `task check` |
| Get the same result as CI (incl. Go race detector on Windows) | `task docker:check` |
| Run the service locally in Docker | `task docker:up` -> http://localhost:8080 (`/healthz`) |
| See logs / stop | `task docker:logs` / `task docker:down` |
| Rebuild automatically while coding | `task docker:watch` |
| See every command | `task --list` |

| Track | Check = | Gotchas to remember |
|---|---|---|
| **Go** | `go vet` + `gofmt -l` + `go test -race` + `go build` | `-race` needs a C compiler (skipped on plain Windows: use `docker:check`) |
| **.NET** | `dotnet format --verify-no-changes` + `dotnet build` (warnings = errors) + `dotnet test` | CRLF breaks `dotnet format`; SDK pinned in `global.json`; run `task fmt` after hand edits |

Windows rules: WSL2 backend; agent and repo in the same environment (all-WSL or all-Windows); LF line endings; install Task with `winget install Task.Task`; confirm Docker Desktop licensing/managed install.

---

## 9. When the agent gets it wrong

1. Was the task small and specific? Split it.
2. Was the needed context in the repo? Add it.
3. Could the agent verify itself? Fix the check.
4. Is the session long or confused? Hand off and restart.
5. Still wrong? Do that part yourself.

Then **write the lesson down**: one line in `AGENTS.md`, or a tweak to a prompt file. That is how the whole team gets better.

---

## 10. Ownership and changes to this playbook

| Item | Owner |
|---|---|
| Playbook and templates | `<engineering leadership / platform team>` |
| Tool approvals and policy | `<security + legal>` |
| Champions per team | `<names>` |
| Review cadence | monthly; propose changes via PR to this file |
