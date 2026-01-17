# Complete Documentation Index

> **Master Index**: All documentation files with descriptions, keywords, and navigation paths.
> **For AI Models**: Use this to quickly locate any information you need.

## 🎯 Start Here

### For AI Models

1. **[PROJECT_CONTEXT.md](PROJECT_CONTEXT.md)** - Read this FIRST for complete context
2. **[AI_GUIDE.md](AI_GUIDE.md)** - AI-specific workflows and analysis patterns
3. **[CONVENTIONS.md](CONVENTIONS.md)** - Coding standards (must follow)
4. **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Fast lookups and snippets

### For Humans

1. **[docs/NEXT_STEPS.md](docs/NEXT_STEPS.md)** - What to do now
2. **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)** - Add features
3. **[docs/CONFIGURATION_GUIDE.md](docs/CONFIGURATION_GUIDE.md)** - Configure CLI

---

## 📚 All Documentation Files

### Root Level Context Files

| File | Type | Purpose | Keywords |
|------|------|---------|----------|
| **PROJECT_CONTEXT.md** | Context | Complete project overview, architecture, structure | architecture, overview, context, structure, layers |
| **AI_GUIDE.md** | Guide | AI-specific workflows and analysis patterns | ai, workflow, analysis, patterns, best practices |
| **CONVENTIONS.md** | Standards | Coding standards, naming, patterns, rules | conventions, standards, patterns, rules, naming |
| **QUICK_REFERENCE.md** | Reference | Quick lookups, snippets, common patterns | reference, snippets, quick, lookup, cheat sheet |
| **README.md** | Overview | Project README, getting started | readme, project, overview, getting started |

### Detailed Guides (docs/)

| File | Type | Purpose | Keywords |
|------|------|---------|----------|
| **docs/README.md** | Index | Documentation index with navigation | index, navigation, documentation, map |
| **docs/NEXT_STEPS.md** | Guide | What to do after setup | next steps, getting started, post-setup, first steps |
| **docs/GETTING_STARTED.md** | Tutorial | How to add commands and features | add command, create feature, tutorial, how-to |
| **docs/CONFIGURATION_GUIDE.md** | Reference | Complete configuration reference | configuration, config, settings, api, database, files, output |
| **docs/CUSTOMIZATION_COMPLETE.md** | Summary | What has been customized | customization, changes, modifications, summary |
| **docs/PROJECT_SUMMARY.md** | Overview | Full project summary | summary, overview, features, structure |
| **docs/GO_CLI_COMPREHENSIVE_REFERENCE.md** | Reference | Authoritative Go CLI patterns | cobra, viper, cli, patterns, best practices, testing |
| **docs/GO_CLI_SETUP_GUIDE.md** | Guide | Original setup instructions | setup, installation, initialization |

### Configuration Files

| File | Type | Purpose | Keywords |
|------|------|---------|----------|
| **configs/config.example.yaml** | Example | Complete configuration example with comments | config, example, yaml, settings, all options |

### Code Documentation

| File | Type | Purpose | Keywords |
|------|------|---------|----------|
| **internal/config/config.go** | Code | Configuration structures and validation | config struct, validation, types |
| **internal/config/defaults.go** | Code | Default configuration values | defaults, default values, viper defaults |
| **internal/model/errors.go** | Code | Domain error types | errors, domain errors, validation error |
| **internal/output/formatter.go** | Code | Output formatting utility | formatter, output, json, yaml, table, csv |
| **internal/logger/logger.go** | Code | Logging configuration | logger, slog, logging setup |
| **cmd/example/greet.go** | Example | Complete command example | example, command, pattern, reference |
| **internal/handler/greet.go** | Example | Handler pattern example | handler, example, pattern |
| **internal/service/example/service.go** | Example | Service pattern example | service, example, business logic |

---

## 🔍 Search by Topic

### Architecture & Design

**Primary**: PROJECT_CONTEXT.md, CONVENTIONS.md
**Secondary**: docs/PROJECT_SUMMARY.md, docs/GO_CLI_COMPREHENSIVE_REFERENCE.md

**Topics**: Clean architecture, layers, dependencies, structure

---

### Adding Features / Commands

**Primary**: docs/GETTING_STARTED.md, QUICK_REFERENCE.md
**Secondary**: cmd/example/greet.go, AI_GUIDE.md

**Topics**: Add command, create feature, extend, new functionality

---

### Configuration

**Primary**: docs/CONFIGURATION_GUIDE.md, internal/config/config.go
**Secondary**: configs/config.example.yaml, internal/config/defaults.go

**Topics**: Config, settings, API, database, files, output formats, viper

---

### Error Handling

**Primary**: CONVENTIONS.md (Error Handling section), internal/model/errors.go
**Secondary**: AI_GUIDE.md (Pitfalls), docs/GO_CLI_COMPREHENSIVE_REFERENCE.md

**Topics**: Errors, wrapping, domain errors, validation

---

### Logging

**Primary**: CONVENTIONS.md (Logging section), internal/logger/logger.go
**Secondary**: AI_GUIDE.md, QUICK_REFERENCE.md

**Topics**: Logging, slog, structured logging, log levels

---

### Output Formatting

**Primary**: internal/output/formatter.go, docs/CONFIGURATION_GUIDE.md
**Secondary**: QUICK_REFERENCE.md, configs/config.example.yaml

**Topics**: Output, format, JSON, YAML, table, CSV, formatter

---

### Testing

**Primary**: docs/GO_CLI_COMPREHENSIVE_REFERENCE.md (Testing Patterns)
**Secondary**: CONVENTIONS.md (Testing), docs/GETTING_STARTED.md

**Topics**: Testing, table-driven tests, mocks, unit tests

---

### CLI Framework (Cobra)

**Primary**: docs/GO_CLI_COMPREHENSIVE_REFERENCE.md (CLI Framework section)
**Secondary**: cmd/example/greet.go, QUICK_REFERENCE.md

**Topics**: Cobra, commands, flags, arguments

---

## 🗺️ Documentation Navigation Map

```
                    START
                      ↓
            ┌─────────┴─────────┐
            │                   │
         AI Model            Human
            │                   │
            ↓                   ↓
    PROJECT_CONTEXT.md    docs/NEXT_STEPS.md
            ↓                   ↓
      AI_GUIDE.md         docs/GETTING_STARTED.md
            ↓                   ↓
    CONVENTIONS.md        docs/CONFIGURATION_GUIDE.md
            ↓                   ↓
  QUICK_REFERENCE.md      Specific docs as needed
            ↓
    Task-specific docs
```

---

## 📋 Documentation by File Count

- **Root Level**: 5 files (PROJECT_CONTEXT, AI_GUIDE, CONVENTIONS, QUICK_REFERENCE, README)
- **docs/**: 8 files (README, 7 guides)
- **configs/**: 1 file (config.example.yaml)
- **Code docs**: 7 files (config, model, output, logger, examples)

**Total**: 21 documentation files

---

## 🎯 Quick Task Lookup

| Task | Start Here | Then Read | Example |
|------|------------|-----------|---------|
| Add command | docs/GETTING_STARTED.md | QUICK_REFERENCE.md | cmd/example/greet.go |
| Configure API | docs/CONFIGURATION_GUIDE.md | internal/config/config.go | configs/config.example.yaml |
| Handle errors | CONVENTIONS.md#error-handling | internal/model/errors.go | - |
| Format output | internal/output/formatter.go | docs/CONFIGURATION_GUIDE.md | - |
| Add tests | docs/GO_CLI_COMPREHENSIVE_REFERENCE.md#6 | CONVENTIONS.md#testing | - |
| Understand architecture | PROJECT_CONTEXT.md | CONVENTIONS.md | docs/PROJECT_SUMMARY.md |

---

## 💡 Documentation Principles

This documentation is organized to be:

1. **Discoverable**: Multiple entry points for different audiences
2. **Hierarchical**: From overview to details
3. **Cross-referenced**: Related docs link to each other
4. **Task-oriented**: Organized by what you want to do
5. **Keyword-rich**: Easy to search and find
6. **Example-driven**: Working examples for all patterns

---

## 🚀 Reading Order Recommendations

### For AI Models Working on a Task

```
1. PROJECT_CONTEXT.md      (5 min - architecture & structure)
   ↓
2. AI_GUIDE.md             (5 min - AI-specific workflows)
   ↓
3. CONVENTIONS.md          (10 min - coding standards)
   ↓
4. Task-specific doc       (varies - depends on task)
   ↓
5. QUICK_REFERENCE.md      (2 min - as needed for snippets)
```

### For Developers Adding a Feature

```
1. docs/NEXT_STEPS.md               (2 min - orientation)
   ↓
2. docs/GETTING_STARTED.md          (10 min - how to add features)
   ↓
3. QUICK_REFERENCE.md               (5 min - code snippets)
   ↓
4. cmd/example/greet.go             (5 min - working example)
   ↓
5. Start coding!
```

### For Understanding the Codebase

```
1. PROJECT_CONTEXT.md               (overview)
   ↓
2. docs/PROJECT_SUMMARY.md          (details)
   ↓
3. CONVENTIONS.md                   (patterns)
   ↓
4. docs/GO_CLI_COMPREHENSIVE_REFERENCE.md (deep dive)
```

---

## 🔗 External Links

- **Main README**: [README.md](README.md)
- **Documentation Hub**: [docs/README.md](docs/README.md)
- **Example Code**: [cmd/example/](cmd/example/)
- **Configuration**: [configs/](configs/)

---

## ✅ Documentation Checklist

When working on this project, have you:

- [ ] Read PROJECT_CONTEXT.md?
- [ ] Checked CONVENTIONS.md for standards?
- [ ] Reviewed existing examples in cmd/example/?
- [ ] Consulted task-specific documentation?
- [ ] Used QUICK_REFERENCE.md for snippets?

If yes to all, you're ready to work effectively!

---

**Last Updated**: 2026-01-18
**Maintained By**: Project team
**Status**: Comprehensive documentation for optimal AI and developer experience
