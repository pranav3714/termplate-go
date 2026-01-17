# Documentation Index

> **For AI Models**: This is the central documentation hub. All guides are organized below by purpose and task type.
> **For Humans**: Start with [Next Steps](#quick-start-for-humans) or [Getting Started](#development-guides).

## Table of Contents

- [Quick Start (For Humans)](#quick-start-for-humans)
- [Quick Start (For AI Models)](#quick-start-for-ai-models)
- [Documentation by Purpose](#documentation-by-purpose)
- [Documentation by Task](#documentation-by-task)
- [All Documents](#all-documents)

---

## Quick Start (For Humans)

**New to the project?** Read these in order:

1. **[Next Steps](NEXT_STEPS.md)** ← Start here
2. **[Getting Started](GETTING_STARTED.md)** - How to add features
3. **[Configuration Guide](CONFIGURATION_GUIDE.md)** - How to configure

**Need something specific?**

- Add a command: [Getting Started - Adding Commands](GETTING_STARTED.md#adding-a-new-command)
- Configure settings: [Configuration Guide](CONFIGURATION_GUIDE.md)
- Understand architecture: [Project Summary - Architecture](PROJECT_SUMMARY.md#architecture-flow)

---

## Quick Start (For AI Models)

**For analysis and development tasks:**

1. **Start with**: `../PROJECT_CONTEXT.md` - Complete project overview and architecture
2. **Check conventions**: `../CONVENTIONS.md` - Coding standards and patterns
3. **Quick lookups**: `../QUICK_REFERENCE.md` - Common snippets and file locations
4. **Specific tasks**: See [Documentation by Task](#documentation-by-task) below

**Architecture at a glance**: cmd → handler → service → repository → model

**Key principle**: Dependencies only flow downward, no circular imports.

---

## Documentation by Purpose

### 📚 Learning & Onboarding

| Document | Purpose | Keywords |
|----------|---------|----------|
| [Next Steps](NEXT_STEPS.md) | What to do after setup | next, start, begin, first steps |
| [Getting Started](GETTING_STARTED.md) | Add features and customize | add command, create feature, extend |
| [Project Summary](PROJECT_SUMMARY.md) | Complete project overview | overview, summary, what's included |

### 🔧 Configuration & Setup

| Document | Purpose | Keywords |
|----------|---------|----------|
| [Configuration Guide](CONFIGURATION_GUIDE.md) | All configuration options | config, settings, api, database, files |
| [Customization Complete](CUSTOMIZATION_COMPLETE.md) | What's been customized | customization, changes, modifications |

### 📖 Reference

| Document | Purpose | Keywords |
|----------|---------|----------|
| [CLI Comprehensive Reference](GO_CLI_COMPREHENSIVE_REFERENCE.md) | Complete CLI patterns | cobra, viper, patterns, best practices |
| [CLI Setup Guide](GO_CLI_SETUP_GUIDE.md) | Original setup guide | setup, installation, initialize |

---

## Documentation by Task

### Task: "I want to add a new command"

1. **Read**: [Getting Started - Adding a New Command](GETTING_STARTED.md#adding-a-new-command)
2. **Reference**: `../QUICK_REFERENCE.md` - Command snippets
3. **Example**: Check `../cmd/example/greet.go`
4. **Pattern**: Follow cmd → handler → service flow

**Files to create**:
- `cmd/mycommand/mycommand.go`
- `internal/handler/myhandler.go`
- `internal/service/myservice/service.go`

**Files to modify**:
- `cmd/root.go` (register command)

---

### Task: "I want to configure API/Database/Files"

1. **Read**: [Configuration Guide](CONFIGURATION_GUIDE.md)
2. **Check**: `../configs/config.example.yaml` - All options
3. **Edit**: `../internal/config/config.go` - Add custom fields
4. **Defaults**: `../internal/config/defaults.go` - Set defaults

**Configuration sections**:
- API: `api.*` (keys, tokens, retry, rate limiting)
- Files: `files.*` (paths, patterns, size limits)
- Database: `database.*` (driver, connection pooling)
- Output: `output.*` (format, colors, table styles)

---

### Task: "I want to understand the architecture"

1. **Read**: `../PROJECT_CONTEXT.md` - Architecture overview
2. **Read**: `../CONVENTIONS.md` - Layer responsibilities
3. **Read**: [Project Summary - Architecture Flow](PROJECT_SUMMARY.md#architecture-flow)

**Key files**:
- `cmd/root.go` - Entry point
- `internal/handler/greet.go` - Handler example
- `internal/service/example/service.go` - Service example

---

### Task: "I want to format output (JSON/YAML/Table/CSV)"

1. **Check**: `../internal/output/formatter.go` - Formatter implementation
2. **Read**: [Configuration Guide - Output Formatting](CONFIGURATION_GUIDE.md#output-formatting)
3. **Use**: `formatter.Print(data)` - Automatic formatting

**Supported formats**: text, json, yaml, table (ASCII/Unicode/Markdown), csv

---

### Task: "I want to handle errors properly"

1. **Check**: `../internal/model/errors.go` - Domain errors
2. **Read**: `../CONVENTIONS.md#error-handling` - Error patterns
3. **Pattern**: Always wrap with `fmt.Errorf("context: %w", err)`

**Domain errors available**:
- `model.ErrNotFound`
- `model.ErrAlreadyExists`
- `model.NewValidationError(field, message)`
- `model.NewOperationError(op, entity, id, err)`

---

### Task: "I want to add tests"

1. **Read**: [Getting Started - Testing](GETTING_STARTED.md#testing)
2. **Read**: [CLI Reference - Testing Patterns](GO_CLI_COMPREHENSIVE_REFERENCE.md#6-testing-patterns)
3. **Pattern**: Table-driven tests

**Example**: `internal/service/example/service_test.go` (create this)

---

### Task: "I want to add logging"

1. **Check**: `../internal/logger/logger.go` - Logger setup
2. **Read**: `../CONVENTIONS.md#logging` - Logging patterns
3. **Use**: `slog.Info("message", "key", value)`

**Levels**: Debug, Info, Warn, Error

---

## All Documents

### Root Level (Project Context)

- **[PROJECT_CONTEXT.md](../PROJECT_CONTEXT.md)** - Complete project context for AI analysis
- **[CONVENTIONS.md](../CONVENTIONS.md)** - Coding standards and patterns
- **[QUICK_REFERENCE.md](../QUICK_REFERENCE.md)** - Quick lookup for common tasks
- **[README.md](../README.md)** - Main project README

### docs/ Directory (Detailed Guides)

#### Getting Started

- **[NEXT_STEPS.md](NEXT_STEPS.md)** - What to do now after setup
- **[GETTING_STARTED.md](GETTING_STARTED.md)** - How to add commands and features
- **[CUSTOMIZATION_COMPLETE.md](CUSTOMIZATION_COMPLETE.md)** - Summary of customizations

#### Configuration

- **[CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md)** - Complete configuration reference

#### Reference

- **[GO_CLI_COMPREHENSIVE_REFERENCE.md](GO_CLI_COMPREHENSIVE_REFERENCE.md)** - Authoritative Go CLI patterns
- **[GO_CLI_SETUP_GUIDE.md](GO_CLI_SETUP_GUIDE.md)** - Original setup guide
- **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Full project summary

#### Command Documentation

- **[commands/](commands/)** - Command reference documentation (empty, ready for use)

---

## Documentation Organization

### By Audience

**For New Users**:
1. Next Steps → Getting Started → Configuration Guide

**For Developers**:
1. PROJECT_CONTEXT.md → CONVENTIONS.md → GETTING_STARTED.md

**For AI Models**:
1. PROJECT_CONTEXT.md → CONVENTIONS.md → QUICK_REFERENCE.md → Specific guides

### By Topic

| Topic | Documents |
|-------|-----------|
| **Architecture** | PROJECT_CONTEXT.md, CONVENTIONS.md, PROJECT_SUMMARY.md |
| **Configuration** | CONFIGURATION_GUIDE.md, config/config.go, configs/config.example.yaml |
| **Development** | GETTING_STARTED.md, CONVENTIONS.md, QUICK_REFERENCE.md |
| **CLI Patterns** | GO_CLI_COMPREHENSIVE_REFERENCE.md, CONVENTIONS.md |
| **Getting Started** | NEXT_STEPS.md, GETTING_STARTED.md, CUSTOMIZATION_COMPLETE.md |

---

## File Structure Overview

```
Documentation locations:

Root level (context & reference):
├── PROJECT_CONTEXT.md          # Project overview & architecture
├── CONVENTIONS.md              # Coding standards
├── QUICK_REFERENCE.md          # Fast lookups
└── README.md                   # Main project README

docs/ (detailed guides):
├── README.md                   # This file (documentation index)
├── NEXT_STEPS.md              # Getting started
├── GETTING_STARTED.md         # Adding features
├── CONFIGURATION_GUIDE.md     # Configuration reference
├── CUSTOMIZATION_COMPLETE.md  # Customization summary
├── PROJECT_SUMMARY.md         # Project overview
├── GO_CLI_COMPREHENSIVE_REFERENCE.md  # CLI patterns
└── GO_CLI_SETUP_GUIDE.md      # Setup guide

Code documentation:
├── cmd/                        # See example in cmd/example/
├── internal/config/config.go   # Configuration structures
├── internal/output/formatter.go # Output formatting
└── configs/config.example.yaml # Configuration template
```

---

## Quick Answers

### "Where do I start?"
→ `NEXT_STEPS.md`

### "How do I add a command?"
→ `GETTING_STARTED.md` section "Adding a New Command"

### "How do I configure the CLI?"
→ `CONFIGURATION_GUIDE.md`

### "What patterns should I follow?"
→ `CONVENTIONS.md` and `GO_CLI_COMPREHENSIVE_REFERENCE.md`

### "What's the architecture?"
→ `PROJECT_CONTEXT.md`

### "Where are the examples?"
→ `../cmd/example/` directory

### "What's the quick reference?"
→ `../QUICK_REFERENCE.md`

---

## Navigation Tips

### For AI Models

When analyzing or working on tasks:

1. **Always start** with `../PROJECT_CONTEXT.md` for full context
2. **Check** `../CONVENTIONS.md` for standards
3. **Use** `../QUICK_REFERENCE.md` for fast lookups
4. **Reference** specific guides in `docs/` for detailed information
5. **Follow** existing examples in `../cmd/example/`

### For Developers

When adding features:

1. **Plan**: Review architecture in `PROJECT_CONTEXT.md`
2. **Learn**: Read `GETTING_STARTED.md`
3. **Reference**: Use `QUICK_REFERENCE.md` for snippets
4. **Configure**: Check `CONFIGURATION_GUIDE.md`
5. **Follow**: Standards in `CONVENTIONS.md`

---

## Document Relationships

```
PROJECT_CONTEXT.md (overview) ──┬── CONVENTIONS.md (how to code)
                                │   └── GO_CLI_COMPREHENSIVE_REFERENCE.md (patterns)
                                │
                                ├── QUICK_REFERENCE.md (fast lookup)
                                │
                                └── docs/
                                    ├── NEXT_STEPS.md (what now?)
                                    ├── GETTING_STARTED.md (add features)
                                    ├── CONFIGURATION_GUIDE.md (settings)
                                    └── PROJECT_SUMMARY.md (full overview)
```

---

**Last Updated**: 2026-01-18
**Total Documents**: 11 (3 root-level + 8 in docs/)
**Status**: Comprehensive documentation for AI and human developers
