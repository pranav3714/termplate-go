# Security Policy

## 🔒 Reporting Security Vulnerabilities

If you discover a security vulnerability in Termplate Go, please report it by:

1. **Opening a GitHub Security Advisory**: [Create Security Advisory](https://github.com/blacksilver/termplate-go/security/advisories/new)
2. **OR emailing** the maintainers (if you prefer private disclosure)

**Please do NOT** open a public issue for security vulnerabilities.

We will respond as quickly as possible to address the issue.

---

## ⚠️ CRITICAL: Never Commit Secrets!

### What NOT to Commit

**NEVER commit any of the following to the repository**:

- ❌ API keys or tokens
- ❌ Passwords or credentials
- ❌ Private keys or certificates
- ❌ Database credentials
- ❌ Environment files (`.env`, `.env.local`, etc.)
- ❌ Service account files
- ❌ OAuth tokens or secrets
- ❌ SSH keys
- ❌ AWS credentials
- ❌ Any file containing sensitive data

### Our .gitignore Protects You

The `.gitignore` file is configured to prevent committing:

```gitignore
# Environment files
.env
.env.*
*.secret
secrets.yaml

# Credentials
*.key
*.pem
credentials.json
*.token

# Config files that may contain secrets
.termplate.yaml
*.local.yaml
*.prod.yaml

# And many more patterns...
```

**However**, `.gitignore` is not foolproof! Always double-check before committing.

---

## ✅ How to Handle Secrets Safely

### 1. Use Environment Variables

**DO THIS**:
```bash
# Set environment variables
export TERMPLATE_API_KEY=your-secret-key
export TERMPLATE_DB_PASSWORD=your-db-password

# Run your CLI
./build/bin/termplate your-command
```

**In config files, use variable substitution**:
```yaml
api:
  key: ${TERMPLATE_API_KEY}  # ✅ GOOD - references env var

database:
  password: ${TERMPLATE_DB_PASSWORD}  # ✅ GOOD
```

### 2. Use Example Files

**DO THIS**:
```bash
# Create example file (safe to commit)
cp config.yaml config.example.yaml

# Add to git
git add config.example.yaml  # ✅ Safe - no real secrets

# Users copy and add their own secrets
cp config.example.yaml ~/.termplate.yaml  # User's local config
```

**In example files, show the pattern**:
```yaml
api:
  key: ${TERMPLATE_API_KEY}  # ✅ GOOD - placeholder
  # or
  key: "your-api-key-here"   # ✅ GOOD - obviously placeholder
```

### 3. Check Before Committing

**Always run these checks**:

```bash
# Check what you're about to commit
git diff --staged

# Look for potential secrets
git diff --staged | grep -i "api.key\|password\|secret\|token"

# If you find any, DON'T COMMIT
git reset HEAD <file>
```

### 4. Use Git Hooks (Already Set Up!)

This project has lefthook configured to run checks before commits. However, **git hooks are not a replacement for manual vigilance**.

---

## 🚨 What If I Accidentally Committed a Secret?

**If you accidentally commit a secret to a public repository, it's already compromised!**

### Immediate Actions

1. **Revoke the secret immediately**
   - Rotate API keys
   - Change passwords
   - Invalidate tokens

2. **Remove from Git history** (if caught early):
   ```bash
   # If it's the last commit and not pushed yet
   git reset --soft HEAD~1
   git reset HEAD <file-with-secret>
   # Remove the secret from the file
   git add <file>
   git commit -m "Your message"
   ```

3. **If already pushed**, the secret is compromised:
   - Revoke it immediately
   - Consider the secret public
   - Don't rely on removing it from history

### Prevention Tools

Consider using:

- **git-secrets**: Prevents committing secrets
- **gitleaks**: Scans for secrets in commits
- **pre-commit hooks**: Additional checks before commits

---

## 🔐 Best Practices for This Project

### For Configuration

1. **Always use environment variables** for secrets:
   ```go
   apiKey := os.Getenv("TERMPLATE_API_KEY")
   if apiKey == "" {
       return fmt.Errorf("TERMPLATE_API_KEY environment variable not set")
   }
   ```

2. **Use Viper's environment variable support**:
   ```yaml
   # config.yaml
   api:
     key: ${TERMPLATE_API_KEY}  # Viper will substitute from env
   ```

3. **Provide clear example files**:
   - `config.example.yaml` ✅ (safe to commit)
   - `.termplate.yaml` ❌ (user's local config with real secrets)

### For Development

1. **Create `.env.example`** with placeholders:
   ```bash
   # .env.example (safe to commit)
   TERMPLATE_API_KEY=your-api-key-here
   TERMPLATE_DB_PASSWORD=your-password-here
   ```

2. **Users copy to `.env`** (gitignored):
   ```bash
   cp .env.example .env
   # Edit .env with real secrets
   ```

3. **Never commit** `.env` or any file with real secrets

### For Code Review

When reviewing pull requests:

- ❌ Reject PRs with hardcoded secrets
- ❌ Reject PRs with real credentials
- ✅ Accept PRs using environment variables
- ✅ Accept PRs with placeholder values in examples

---

## 📋 Security Checklist for Contributors

Before submitting a PR, verify:

- [ ] No hardcoded API keys, passwords, or tokens
- [ ] All secrets use environment variables
- [ ] Example configs use placeholders only
- [ ] No credentials in test files
- [ ] No real database credentials
- [ ] Ran `git diff --staged` to review changes
- [ ] Checked for accidental secret inclusion

---

## 🛡️ Additional Security Considerations

### Input Validation

- Always validate user input
- Sanitize file paths to prevent directory traversal
- Validate URLs before making requests
- Check file sizes before processing

### Dependencies

- Keep dependencies up to date
- Run `make vuln` regularly to check for vulnerabilities
- Review dependency changes in PRs

### Logging

- **Never log secrets** (API keys, passwords, tokens)
- **Never log sensitive user data** (emails, personal info)
- Sanitize logs before outputting

```go
// ❌ BAD - logs the API key
slog.Info("making request", "api_key", apiKey)

// ✅ GOOD - doesn't log secrets
slog.Info("making request", "url", url)
```

---

## 📞 Contact

For security concerns or questions:

- **Security vulnerabilities**: Use GitHub Security Advisories
- **Questions about secrets**: Open a discussion or issue
- **Need help**: Check [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 🎯 Remember

**The `.gitignore` file helps, but YOU are the last line of defense.**

Always think before you commit:
- "Does this file contain secrets?"
- "Would I be okay with this being public?"
- "Am I committing example data or real data?"

When in doubt, don't commit it! 🔒

---

**Last Updated**: 2026-01-18
