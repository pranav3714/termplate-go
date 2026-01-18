# Changelog

All notable changes to Termplate Go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Renamed project from "ever-so-powerful-go" to "Termplate Go"
- Updated all documentation with new project name
- Enhanced .gitignore to prevent secret commits

### Added
- Comprehensive GitHub compliance files
- Enhanced security documentation
- Issue templates for bugs and features
- Pull request template

## [0.1.0] - 2026-01-18

### Added

**Foundation**:
- Complete Go CLI structure using Cobra framework
- Viper configuration system with multiple sources
- Structured logging with slog (JSON/text output)
- Clean architecture (cmd → handler → service → repository → model)
- Working example command demonstrating full architecture

**Configuration**:
- API client configuration (keys, tokens, retries, rate limiting, headers)
- File processing configuration (patterns, size limits, permissions, backup)
- Database configuration (PostgreSQL, MySQL, SQLite with connection pooling)
- Output formatting configuration (text, JSON, YAML, table, CSV)
- Environment variable support with `TERMPLATE_` prefix
- Configuration validation and helper methods

**Output Formatting**:
- Multi-format output formatter (JSON, YAML, Table, CSV)
- Three table styles (ASCII, Unicode, Markdown)
- Pretty printing support
- Color output control

**Development Tools**:
- Makefile with 16 targets (build, test, lint, audit, etc.)
- golangci-lint with 20+ linters configured
- lefthook for git hooks (pre-commit: format/lint, pre-push: test/vet)
- GitHub Actions workflows (CI and Release)
- GoReleaser configuration for multi-platform releases
- Docker support with multi-stage builds

**Documentation** (21 files, 7,700+ lines):
- PROJECT_CONTEXT.md - Complete project overview
- AI_GUIDE.md - AI model workflows and analysis patterns
- CONVENTIONS.md - Coding standards and patterns
- QUICK_REFERENCE.md - Fast lookups and snippets
- DOCUMENTATION_INDEX.md - Master documentation index
- Comprehensive guides in docs/ directory
- Configuration guide with all options
- Getting started guide with examples

**Quality Assurance**:
- Comprehensive .gitignore (prevents secret commits)
- Error handling with domain-specific errors
- Table-driven test patterns
- Code quality enforcement with linters

### Features

- Shell completion for Bash, Zsh, Fish, PowerShell
- Version command with multiple output formats
- Example command showing full architecture
- Context-aware signal handling (SIGINT, SIGTERM)
- Cross-platform builds (Linux, macOS, Windows - amd64, arm64)

### Documentation

- AI-optimized documentation structure
- Multiple navigation paths (by audience, task, topic)
- Metadata headers on key documentation
- Task-oriented guides
- Working code examples
- Honest disclaimer about being a weekend project 😄

---

## Version History

### Version Format

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backwards-compatible)
- **PATCH**: Bug fixes (backwards-compatible)

### Release Process

1. Update CHANGELOG.md with version and date
2. Create git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
3. Push tag: `git push origin v1.0.0`
4. GitHub Actions will automatically build and release

---

## Links

- [Unreleased changes](https://github.com/pranav3714/termplate-go/compare/v0.1.0...HEAD)
- [All releases](https://github.com/pranav3714/termplate-go/releases)

[Unreleased]: https://github.com/pranav3714/termplate-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/pranav3714/termplate-go/releases/tag/v0.1.0
