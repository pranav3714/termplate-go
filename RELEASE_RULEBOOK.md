# Release Rulebook

> **Audience**: Developers, maintainers, AI assistants
> **Keywords**: release, deploy, tag, version, publish, changelog, production, shipping, rollout
> **Last Updated**: 2026-01-18

Complete guide for creating releases in Termplate Go. This document covers automated release processes, troubleshooting, and best practices.

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Automated Process](#automated-process)
5. [Manual Process](#manual-process)
6. [Automation Script Reference](#automation-script-reference)
7. [AI-Assisted Releases](#ai-assisted-releases)
8. [Troubleshooting](#troubleshooting)
9. [Rollback Procedures](#rollback-procedures)
10. [Version Numbering](#version-numbering)
11. [Security Considerations](#security-considerations)
12. [Quick Reference](#quick-reference)

---

## Overview

### Release System

Termplate Go uses a **tag-triggered release system** powered by:
- **GitHub Actions** - Automated CI/CD pipeline
- **GoReleaser** - Multi-platform binary builds
- **Keep a Changelog** - Structured changelog format
- **Semantic Versioning** - Version numbering scheme

### Current Version

- **Latest**: v0.1.1 (2026-01-18)
- **Repository**: https://github.com/pranav3714/termplate-go
- **Releases**: https://github.com/pranav3714/termplate-go/releases

### Release Trigger

Releases are triggered automatically when a git tag matching `v*` is pushed to the repository:

```bash
git push origin v0.2.0  # Triggers release automation
```

### Build Artifacts

Each release automatically builds for:
- **Linux**: amd64, arm64
- **macOS**: amd64, arm64
- **Windows**: amd64

Archives include: `LICENSE`, `README.md`, `configs/config.example.yaml`

---

## Prerequisites

Before creating a release, ensure:

### Required

- ✅ **Clean git working directory** - No uncommitted changes
- ✅ **All changes committed** - Work in progress must be committed
- ✅ **Tests passing** - Run `make test`
- ✅ **Linters passing** - Run `make lint`
- ✅ **On main branch** - Or intentionally on another branch
- ✅ **Synced with remote** - Pull latest changes

### Recommended

- 📝 **CHANGELOG.md updated** - Document changes in [Unreleased] section
- 🔍 **Code reviewed** - All changes reviewed and approved
- 🧪 **Manual testing completed** - Features tested end-to-end
- 📚 **Documentation updated** - README and docs reflect changes

### Validation Commands

```bash
# Check git status
git status

# Run tests
make test

# Run linters
make lint

# Check for uncommitted changes
git diff
git diff --cached
```

---

## Quick Start

### One-Command Release

```bash
# Interactive release (recommended for first-time use)
make release-prepare
```

Follow the prompts to:
1. Select version number (patch/minor/major or custom)
2. Confirm release details
3. Script automatically handles CHANGELOG, tags, and push

### Preview First

```bash
# Dry-run to preview changes
make release-dry-run
```

This shows what would be changed without making any modifications.

---

## Automated Process

The automated release process uses `scripts/release.sh` to handle all manual steps.

### Step 1: Prepare Changes

Ensure all code changes are committed and pushed:

```bash
git add .
git commit -m "feat: your feature"
git push origin main
```

### Step 2: Update CHANGELOG (Optional)

If [Unreleased] section is empty, add your changes:

```markdown
## [Unreleased]

### Added
- New feature description

### Changed
- Modified behavior description

### Fixed
- Bug fix description
```

### Step 3: Run Release Script

```bash
# Interactive mode (recommended)
make release-prepare

# Or with automatic version increment
make release-patch  # 0.1.1 → 0.1.2
```

### Step 4: Automated Actions

The script automatically:

1. **Validates prerequisites**
   - Checks git status
   - Runs tests and linters
   - Verifies version format

2. **Updates CHANGELOG.md**
   - Creates backup (CHANGELOG.md.backup)
   - Moves [Unreleased] content to new version section
   - Adds version number and current date
   - Creates empty [Unreleased] section

3. **Updates version links**
   - Updates [Unreleased] comparison link
   - Adds new version comparison link

4. **Creates git commit**
   - Commits CHANGELOG.md changes
   - Message: "chore: prepare release vX.Y.Z"

5. **Creates annotated git tag**
   - Tag name: vX.Y.Z
   - Includes release notes from CHANGELOG
   - Co-authored by Claude

6. **Pushes to remote**
   - Pushes commit to main branch
   - Pushes tag (triggers GitHub Actions)

### Step 5: GitHub Actions

Once the tag is pushed, GitHub Actions automatically:

1. **Builds binaries** for all platforms (6 combinations)
2. **Creates archives** (tar.gz for Unix, zip for Windows)
3. **Generates checksums** (checksums.txt)
4. **Creates GitHub release** with auto-generated changelog
5. **Uploads artifacts** for public download

### Step 6: Verification

Check the release:

```bash
# View local tags
git tag -l

# View remote tags
git ls-remote --tags origin

# Monitor GitHub Actions
open https://github.com/pranav3714/termplate-go/actions

# View releases
open https://github.com/pranav3714/termplate-go/releases
```

---

## Manual Process

If the automation script fails, you can create a release manually.

### Step 1: Update CHANGELOG.md

Edit `CHANGELOG.md`:

```markdown
## [Unreleased]

## [0.2.0] - 2026-01-18

### Added
- Your changes here
```

Update version links at the bottom:

```markdown
[Unreleased]: https://github.com/pranav3714/termplate-go/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/pranav3714/termplate-go/compare/v0.1.1...v0.2.0
```

### Step 2: Commit Changes

```bash
git add CHANGELOG.md
git commit -m "chore: prepare release v0.2.0"
```

### Step 3: Create Tag

```bash
git tag -a v0.2.0 -m "Release v0.2.0

### Added
- Feature 1
- Feature 2

### Changed
- Improvement 1

Co-Authored-By: Developer Name <email@example.com>"
```

### Step 4: Push Tag

```bash
git push origin main
git push origin v0.2.0
```

### Step 5: Monitor Release

GitHub Actions will automatically build and publish the release.

---

## Automation Script Reference

### Command-Line Options

```bash
./scripts/release.sh [OPTIONS]
```

### Options

| Option | Description | Example |
|--------|-------------|---------|
| `--version VERSION` | Specify exact version | `--version v0.2.0` |
| `--patch` | Auto-increment patch (0.1.1→0.1.2) | `--patch` |
| `--minor` | Auto-increment minor (0.1.1→0.2.0) | `--minor` |
| `--major` | Auto-increment major (0.1.1→1.0.0) | `--major` |
| `--dry-run` | Preview without changes | `--dry-run` |
| `--skip-tests` | Skip running tests | `--skip-tests` |
| `--skip-lint` | Skip running linters | `--skip-lint` |
| `--yes` | Auto-confirm prompts | `--yes` |
| `--help` | Show help message | `--help` |

### Usage Examples

```bash
# Interactive mode (prompts for version)
./scripts/release.sh

# Dry-run mode (preview only)
./scripts/release.sh --dry-run

# Specify version explicitly
./scripts/release.sh --version v0.2.0

# Auto-increment patch version
./scripts/release.sh --patch

# Auto-increment minor version
./scripts/release.sh --minor

# Skip validation (not recommended)
./scripts/release.sh --skip-tests --skip-lint

# Fully automated (CI-friendly)
./scripts/release.sh --version v0.2.0 --yes

# Combine options
./scripts/release.sh --patch --dry-run
```

### Makefile Shortcuts

```bash
# Interactive release
make release-prepare

# Preview changes
make release-dry-run

# Auto-increment patch
make release-patch
```

---

## AI-Assisted Releases

### Prompts That Trigger Releases

AI assistants should recognize these prompts:

**Direct commands**:
- "Create a new release"
- "Release version 0.2.0"
- "Prepare release 0.2.0"
- "Ship version 0.2.0"
- "Tag and release 0.2.0"
- "Publish version 0.2.0"

**Incremental**:
- "Release a patch version"
- "Bump patch and release"
- "Create a minor release"
- "Release a new major version"

**Preparation**:
- "Prepare for release"
- "Update changelog for release"
- "Get ready to release"

**Process questions**:
- "How do I release this project?"
- "Show me the release process"
- "Explain the release workflow"
- "What are the release steps?"

### Expected AI Workflow

When asked to create a release, AI assistants should:

1. **Read this document** (RELEASE_RULEBOOK.md)
2. **Validate prerequisites**
   ```bash
   git status
   make test
   make lint
   ```
3. **Execute release**
   ```bash
   make release-prepare
   # or
   ./scripts/release.sh --version v0.2.0
   ```
4. **Monitor output** and answer any prompts
5. **Verify success**
   - Check GitHub Actions status
   - Confirm release appears on GitHub
6. **Report results** with links to:
   - GitHub Actions workflow run
   - Published release page

### Example AI Interaction

```
User: "Release version 0.2.0"

AI Response:
I'll create release v0.2.0 for Termplate Go.

First, let me validate prerequisites:
- Git status: clean ✓
- Tests: passing ✓
- Linters: passing ✓

Executing release:
$ make release-prepare

[Output from script...]

Release v0.2.0 created successfully!

Verification:
- GitHub Actions: https://github.com/pranav3714/termplate-go/actions/runs/...
- Release page: https://github.com/pranav3714/termplate-go/releases/tag/v0.2.0

The release will be available in a few minutes once the build completes.
```

---

## Troubleshooting

### "Working directory not clean"

**Problem**: Uncommitted changes in working directory

**Solution**:
```bash
# Check status
git status

# Commit changes
git add .
git commit -m "your message"

# Or stash changes
git stash
```

### "Tag already exists"

**Problem**: Version tag already created locally or remotely

**Solution**:
```bash
# Check existing tags
git tag -l

# Delete local tag
git tag -d v0.2.0

# Delete remote tag (⚠️ use with caution)
git push origin :refs/tags/v0.2.0

# Then create release with new version
```

### "Tests failing"

**Problem**: Test suite not passing

**Solution**:
```bash
# Run tests to see failures
make test

# Fix failing tests
# Then re-run release

# Or skip tests (not recommended)
./scripts/release.sh --skip-tests
```

### "Linters failing"

**Problem**: Code quality checks not passing

**Solution**:
```bash
# Run linters to see issues
make lint

# Auto-fix issues
make lint-fix

# Then re-run release

# Or skip linters (not recommended)
./scripts/release.sh --skip-lint
```

### "GitHub Actions build failed"

**Problem**: Release workflow failed after tag was pushed

**Solution**:
1. Check GitHub Actions logs:
   ```
   https://github.com/pranav3714/termplate-go/actions
   ```
2. Common issues:
   - GoReleaser configuration error
   - Missing GITHUB_TOKEN permissions
   - Build compilation errors
   - Platform-specific build failures
3. Fix issue and create new patch release

### "Network failure during push"

**Problem**: Failed to push tag to remote

**Solution**:
```bash
# Tag was created locally but not pushed
# Push manually:
git push origin main
git push origin v0.2.0
```

### "Empty [Unreleased] section"

**Problem**: No changes documented in CHANGELOG.md

**Solution**:
- Add release notes manually to CHANGELOG.md [Unreleased] section
- Or continue anyway (script will warn but allow it)

### "Script failed mid-execution"

**Problem**: Script encountered an error after making changes

**Solution**:
```bash
# Restore from backup
cp CHANGELOG.md.backup CHANGELOG.md

# Reset commit if created
git reset HEAD~1

# Delete tag if created
git tag -d v0.2.0

# Fix issue and retry
```

---

## Rollback Procedures

### Before Push (Local Only)

If the release hasn't been pushed yet:

```bash
# Delete local tag
git tag -d v0.2.0

# Restore CHANGELOG from backup
cp CHANGELOG.md.backup CHANGELOG.md

# Reset commit
git reset HEAD~1

# Verify rollback
git status
git tag -l
```

### After Push (Remote)

⚠️ **Warning**: Only rollback if absolutely necessary. Users may have already downloaded the release.

```bash
# Delete remote tag
git push origin :refs/tags/v0.2.0

# Delete local tag
git tag -d v0.2.0

# Manually delete GitHub release via web interface:
# https://github.com/pranav3714/termplate-go/releases
# Click release → Delete

# Fix issues
# Create new patch release
./scripts/release.sh --patch
```

### Best Practices

1. **Never delete published releases** that users have downloaded
2. **Create patch releases** to fix issues instead of rolling back
3. **Use --dry-run** to catch issues before actual release
4. **Test thoroughly** before releasing
5. **Communicate changes** if rollback is necessary

---

## Version Numbering

### Semantic Versioning

Termplate Go follows [Semantic Versioning 2.0.0](https://semver.org/):

**Format**: `vMAJOR.MINOR.PATCH[-PRERELEASE]`

**Examples**: `v1.0.0`, `v0.2.1`, `v2.1.0-beta.1`

### Version Components

```
v1.2.3-beta.1
│ │ │ └─────── Pre-release identifier (optional)
│ │ └───────── PATCH: Bug fixes (backwards-compatible)
│ └─────────── MINOR: New features (backwards-compatible)
└───────────── MAJOR: Breaking changes (incompatible API)
```

### When to Increment

**MAJOR version** (1.0.0 → 2.0.0):
- Remove commands or options
- Change command behavior incompatibly
- Remove configuration options
- Rename flags or commands

**MINOR version** (1.0.0 → 1.1.0):
- Add new commands
- Add new features
- Add new configuration options
- Deprecate features (but keep them working)

**PATCH version** (1.0.0 → 1.0.1):
- Fix bugs
- Fix security vulnerabilities
- Improve documentation
- Internal refactoring
- Performance improvements

### Pre-1.0 Versions

Current status: **v0.1.1** (pre-1.0)

Before 1.0.0:
- Breaking changes allowed in **minor** versions
- Use 0.x.x for development releases
- Release 1.0.0 when API is stable

### Pre-release Versions

For beta/alpha releases:

```bash
v1.0.0-alpha.1    # Alpha release
v1.0.0-beta.1     # Beta release
v1.0.0-rc.1       # Release candidate
```

Usage:
```bash
./scripts/release.sh --version v0.2.0-beta.1
```

---

## Security Considerations

### Secret Protection

**Never include secrets in releases:**
- ❌ `.env` files
- ❌ API keys or tokens
- ❌ Private keys or certificates
- ❌ Database credentials
- ❌ `credentials.json` or similar

### Included Files

Archives only include:
- ✅ `LICENSE`
- ✅ `README.md`
- ✅ `configs/config.example.yaml`
- ✅ Binary executables

### Best Practices

1. **Review .gitignore** - Ensure secrets are ignored
2. **Check archives** - Verify no sensitive data included
3. **Use .goreleaser.yml** - Control what gets packaged
4. **Audit commits** - Review history before releasing
5. **Enable 2FA** - Require on GitHub account
6. **GPG sign tags** (optional) - Add cryptographic signature

### GPG Tag Signing (Optional)

Enable GPG signing for tags:

```bash
# Configure git to sign tags
git config tag.gpgSign true

# Create signed tag
git tag -s v0.2.0 -m "Release v0.2.0"

# Verify signature
git tag -v v0.2.0
```

### Vulnerability Scanning

Before each release:

```bash
# Check for known vulnerabilities
make vuln

# or
govulncheck ./...
```

---

## Quick Reference

### One-Line Commands

```bash
# Interactive release
make release-prepare

# Dry-run preview
make release-dry-run

# Patch release
make release-patch

# Emergency rollback
git tag -d v0.2.0 && git push origin :refs/tags/v0.2.0
```

### Check Release Status

```bash
# List local tags
git tag -l

# List remote tags
git ls-remote --tags origin

# View GitHub Actions
open https://github.com/pranav3714/termplate-go/actions

# View releases
open https://github.com/pranav3714/termplate-go/releases
```

### Validate Before Release

```bash
# Full validation
make audit

# Individual checks
make test
make lint
make vuln
git status
```

### Common Workflows

**Feature release**:
```bash
# After feature development
git add .
git commit -m "feat: new feature"
git push

# Update CHANGELOG [Unreleased] section
# Then:
make release-prepare  # Select 'minor'
```

**Bug fix release**:
```bash
# After bug fix
git add .
git commit -m "fix: bug description"
git push

# Update CHANGELOG [Unreleased] section
# Then:
make release-patch  # Auto-increment patch
```

**Emergency hotfix**:
```bash
# Fix critical bug
git add .
git commit -m "fix: critical issue"
git push

# Quick release
./scripts/release.sh --patch --skip-tests --yes
```

---

## Related Documentation

- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes
- **[CONVENTIONS.md](CONVENTIONS.md)** - Code and commit conventions
- **[AI_GUIDE.md](AI_GUIDE.md)** - AI assistant workflows
- **[.goreleaser.yml](.goreleaser.yml)** - GoReleaser configuration
- **[.github/workflows/release.yml](.github/workflows/release.yml)** - GitHub Actions workflow

---

## Support

### Issues

Report problems with the release process:
- **GitHub Issues**: https://github.com/pranav3714/termplate-go/issues

### Questions

Ask questions about releasing:
- Open a discussion on GitHub
- Reference this document
- Tag with `release` label

---

*Last updated: 2026-01-18*
*Termplate Go Release Rulebook*
