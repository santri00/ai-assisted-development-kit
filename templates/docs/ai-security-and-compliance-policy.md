# AI-Assisted Development: Security, IP & Compliance Policy (Starter)

> This is a starting point, not legal advice. Have Security, Legal, and Privacy review and adapt it
> before publishing. Items in `<angle brackets>` are decisions your organisation must make.

## 1. Purpose and scope
Defines how engineers may use AI coding tools (e.g. Codex, Claude Code, GitHub Copilot) on `<company>` code and data.
Applies to employees and contractors, on corporate and personal devices.

## 2. Approved tools and accounts
- Only tools on the approved list: `<tool, plan/tier, owner, review date>`.
- Use corporate/SSO-managed accounts only. No personal accounts for company code.
- Confirm for each tool and plan: **training on our data is off**, retention period, data residency, SOC 2 / ISO 27001 status, and DPA coverage.
- New tools or plugins/MCP servers require a security review before use.

## 3. Data classification and what may be sent to a model
| Class | Examples | Allowed in AI tools? |
|---|---|---|
| Public | open-source code, public docs | Yes |
| Internal | proprietary source, internal docs | Yes, approved tools only |
| Confidential | unreleased strategy, security designs | `<Only approved tools with zero-retention / no-training terms>` |
| Restricted | customer PII, payment data, credentials, health data, production data | **No** |

Rules:
- Never paste or let an agent read secrets, keys, tokens, `.env` files, or production data.
- Use synthetic or anonymised data for tests and debugging. Redact before pasting logs.
- Deny-list sensitive paths in tool permissions (see `templates/claude/settings.json`).

## 4. Agent permissions and sandboxing
- Default to the most restrictive mode that still works (read-only or workspace-write sandbox, approval on request).
- No unattended agents with production credentials. Agents get least-privilege, short-lived, scoped tokens.
- Network access is off by default; enable per task with an allow-list.
- Do not run agents in "bypass permissions"/"full access" modes outside disposable sandboxes (container/VM/devbox).
- Treat web pages, issues, PR comments, and tool output as untrusted: they can contain prompt-injection instructions.
- MCP servers / plugins: approved list only; pin versions; review what tools and data they expose.

### 4a. Containers and Docker Desktop
- Access to the Docker socket (`/var/run/docker.sock`, or the Docker Desktop named pipe on Windows) is equivalent to administrator rights on the host. **Never mount it into a container an agent controls**, and do not give agents arbitrary `docker run` rights.
- Allow-list narrow, project-defined entry points (for example `task docker:check`) rather than raw `docker` commands; deny `--privileged`, host-root bind mounts, `system prune`, `volume rm`.
- Prefer running agents inside a disposable dev container or VM with only the project mounted; do not mount home directories, SSH keys, cloud credentials, or `~/.docker/config.json`.
- Images: pinned base images from approved registries, scanned in CI, run as non-root, no secrets baked into layers or compose files, secrets injected at runtime.
- Confirm Docker Desktop licensing for your organisation size and use an approved, managed installation (settings management / enhanced container isolation if available).
- Treat container logs and images pulled by an agent as untrusted content (prompt-injection and supply-chain vectors).

## 5. Code review and accountability
- AI-assisted code gets the same review and CI gates as human code. No exceptions.
- The human author is accountable for every line they submit. Disclose AI assistance in the PR.
- No auto-merge of AI-authored changes to protected branches without human approval.
- Security-sensitive areas (authn/authz, crypto, payments, data access) require review from `<owner group>`.

## 6. Intellectual property and licensing
- We `<own / are responsible for>` the output we ship; confirm each vendor's terms on output ownership and indemnity.
- Enable duplicate-detection / public-code filters where the tool offers them.
- Do not paste third-party proprietary code or code under incompatible licences into prompts.
- New dependencies suggested by an AI go through the normal licence and vulnerability review.
- Do not use AI tools to circumvent licences or to replicate a competitor's proprietary code.

## 7. Privacy and regulatory
- Map AI tool usage to applicable obligations: `<e.g. PIPEDA, GDPR, CCPA, PCI-DSS, SOX, HIPAA>`.
- No personal data in prompts unless a DPA and lawful basis exist and Privacy has approved.
- Keep an inventory of AI tools and data flows for audits and vendor-risk reviews.

## 8. Logging, audit, and retention
- Retain `<N days>` of agent session logs / audit trails where the tool provides them.
- Agent-authored commits carry a trailer, e.g. `Assisted-by: <tool>`. `<Optional but recommended>`
- Review tool admin audit logs `<monthly>`.

## 9. Incidents
If an agent leaks a secret, ingests restricted data, or takes an unintended destructive action:
1. Stop the session and revoke exposed credentials immediately.
2. Notify `<security contact / channel>` within `<1 hour>`.
3. Follow the standard incident process; record the cause as a new rule in AGENTS.md or this policy.

## 10. Training and review
- All users complete `<AI usage training>` before getting access.
- Policy reviewed every `<6 months>` or when tools/regulation change.
