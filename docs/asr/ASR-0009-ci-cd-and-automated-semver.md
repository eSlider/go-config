# ASR-0009: CI/CD pipeline and automated semantic versioning

## Context

ASR-0008 only covered GitHub presentation and the initial manual tag of `v0.1.0`.
The project needs continuous verification on every PR and a release flow that does
not require humans to bump version strings. The source tree must stay free of
version literals (the CHANGELOG aside) so that `go install …@latest` and GoReleaser
artefacts agree on a single source of truth: **git tags**.

## Decision

1. **Test pipeline** (`.github/workflows/test.yml`): matrix of Go `1.22`, `1.23`,
   and `stable` on `ubuntu-latest` plus one `stable` run on macOS and Windows,
   with `go test -race -shuffle=on -count=1 -coverprofile=coverage.out`. Coverage
   is uploaded to Codecov from the Linux `stable` leg only. `go mod tidy` is
   enforced via `git diff --exit-code`.
2. **Lint pipeline** (`.github/workflows/lint.yml`): `go vet` plus
   `golangci-lint@v6` using the repo's existing `.golangci.yml`.
3. **Automated SemVer** (`.github/workflows/release-please.yml` +
   `release-please-config.json` + `.release-please-manifest.json`): Google's
   [release-please-action](https://github.com/googleapis/release-please-action)
   in `release-type: go` mode parses [Conventional
   Commits](https://www.conventionalcommits.org/) since the previous tag and
   opens a "chore: release vX.Y.Z" PR that edits `CHANGELOG.md` and the
   manifest only. Merging the PR creates the git tag. The manifest is seeded
   with the current `0.2.0` so history is continuous.
4. **Release pipeline** (`.github/workflows/release.yml` + `.goreleaser.yml`):
   on `push: tags: ["v*"]`, GoReleaser v2 cross-compiles `cmd/envc` for
   linux/darwin/windows × amd64/arm64 (no windows/arm64), stamps
   `-X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}`
   into the binary, produces `tar.gz` / `zip` archives plus `checksums.txt`,
   and attaches them to a GitHub Release whose body is generated from commit
   groups (`feat`, `fix`, `perf`, …).
5. **No version literals in Go code.** `cmd/envc/main.go` defines
   `var version = "dev"` (and `commit`, `date`), overridable only at link time.
   `envc version` prints whatever the release build injected; local
   `go install` builds report `dev`.
6. **Dependency hygiene** (`.github/dependabot.yml`): weekly `gomod` and
   `github-actions` updates, grouped for minor/patch, with conventional
   `chore(deps)` / `ci(actions)` commit prefixes so release-please classifies
   them correctly.

## Consequences

- Contributors must use Conventional Commits; otherwise release-please will not
  bump the version. `chore:` / `docs:` / `ci:` commits produce no release.
- Breaking changes require `feat!:` or a `BREAKING CHANGE:` footer, which
  release-please maps to a major bump (subject to `bump-minor-pre-major: true`
  while the module is pre-1.0).
- Re-tagging a release manually is discouraged — the tag is the contract that
  triggers GoReleaser.
- Codecov token (`CODECOV_TOKEN`) and branch protection enforcing the Tests
  and Lint checks must be configured in the GitHub repository settings
  (operational, not code).

## Status

Accepted — 2026-05-02.

## References

- [release-please-action](https://github.com/googleapis/release-please-action)
- [GoReleaser v2 docs](https://goreleaser.com/customization/)
- [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
