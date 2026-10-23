# Prompt: Debug a failure

**When:** failing test, production bug, flaky behaviour.

```
Help me debug this problem.

Symptom: <WHAT HAPPENS vs WHAT SHOULD HAPPEN>
Evidence: <ERROR OUTPUT / LOG LINES / FAILING TEST NAME - paste, redact secrets and PII>
Repro: <COMMAND OR STEPS>

Process:
1. Reproduce it first by running the repro. Confirm you see the same failure.
2. List 3 plausible hypotheses ranked by likelihood, and the cheapest experiment to test each.
3. Run the experiments (add temporary logging only; revert it afterwards). Show the evidence.
4. Identify the root cause - not just the symptom - and explain it in plain language.
5. Write a regression test that fails now, then the minimal fix, then run `task check`.

Do not apply speculative fixes. Do not change unrelated code.
```
