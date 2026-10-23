# Prompt: Add tests to existing code (characterization tests)

**When:** legacy or under-tested code, before refactoring.

```
Add tests for <PACKAGE/MODULE/FILE>. Do not change production code.

1. Read the code and existing tests; match the existing test style and helpers.
2. List the behaviours and edge cases you will cover (happy path, boundaries, error paths,
   concurrency if relevant). Show me the list first.
3. Write the tests. Prefer table-driven / parameterised tests. No network, no sleeps, no real clocks.
4. Run them. If a test reveals a likely bug, do NOT fix it silently - report it and mark the test
   clearly (e.g. skipped with a reason) so I can decide.
5. Report coverage change if the toolchain supports it.

Tests must assert behaviour that a regression would break. Avoid tautologies and over-mocking.
```
