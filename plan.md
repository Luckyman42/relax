# Plan: Polishing the `relax` module

Date: 2026-05-16

Goal
----
Make `relax` a polished, trustable open-source Go module that other Go developers will feel comfortable trying and adopting.

Scope
-----
Documentation, tests, CI, code quality, repository/community files, release automation, and packaging hygiene.

Prioritized tasks (High → Low)
--------------------------------

1) Documentation (High)
   - Update `README.md` with correct import path, Quick Start, supported Go versions, and visible badges (CI, coverage, pkg.go.dev, license).
   - Extract runnable examples into an `examples/` directory and reference them from the README.
   - Add package-level docs (`doc.go`) and ensure every exported symbol has a godoc comment.
   - Acceptance: README contains copy-paste runnable examples; pkg.go.dev shows clear docs and examples.

2) Tests & quality (High)
   - Add tests for edge cases: nil pointers, conversions, concurrency, and error-wrapping behavior.
   - Add small benchmarks where relevant and aim for reasonable coverage (goal: improve coverage; realistic target depends on code).
   - Acceptance: tests run in CI with `-race` and report coverage; new tests cover known edge cases.

3) CI & linting (High)
   - Add GitHub Actions workflow to run `go test` across supported Go versions (e.g., 1.21→1.25), `go vet`, and linters (`golangci-lint`).
   - Add a `.golangci.yml` and Make targets for `make test`, `make lint`, `make tidy`.
   - Acceptance: CI workflow passes; badges can be added to README.

4) Repository & community files (Medium)
   - Add `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md` (Contributor Covenant), `SECURITY.md`, issue and PR templates, and optional `CODEOWNERS`.
   - Add `CHANGELOG.md` (Keep a Changelog style).
   - Acceptance: community files present and referenced in README/CONTRIBUTING.

5) Releases & versioning (Medium)
   - Document semantic versioning policy and plan for a `/v2` module path if breaking changes are expected.
   - Add automated releases (goreleaser or GitHub Actions) to produce GitHub releases and artifacts.
   - Acceptance: automated release workflow creates tagged releases and changelogs.

6) Packaging & module hygiene (Medium)
   - Run `go mod tidy` and enforce tidiness in CI.
   - Confirm `module` path in `go.mod` matches repository; add major-version suffix when needed.
   - Acceptance: `go mod tidy` produces no diffs; module path is correct.

7) Polish & usability (Low)
   - Split `relax.go` into logical files (`failer.go`, `guard.go`, `helpers.go`) for readability.
   - Add `Makefile` targets and optional pre-commit hooks for format/lint.
   - Add README migration notes / FAQ and a short