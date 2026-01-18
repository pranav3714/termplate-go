# Pull Request

## Description

<!-- Provide a brief description of your changes -->

## Type of Change

<!-- Mark the relevant option with an 'x' -->

- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Performance improvement
- [ ] Test additions/improvements

## Related Issues

<!-- Link related issues using keywords: Fixes #123, Closes #456, Related to #789 -->

Fixes #

## Changes Made

<!-- List the main changes in this PR -->

-
-
-

## Testing

<!-- Describe how you tested your changes -->

- [ ] Added new tests
- [ ] All tests pass (`make test`)
- [ ] Manual testing completed
- [ ] Tested on multiple platforms (if applicable)

**Test commands used**:
```bash
make test
./build/bin/termplate <your-test-command>
```

## Code Quality

<!-- Verify all quality checks pass -->

- [ ] Code follows project conventions (see [CONVENTIONS.md](../CONVENTIONS.md))
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] No new warnings from `go vet`
- [ ] No vulnerabilities (`make vuln`)
- [ ] All quality checks pass (`make audit`)

## Documentation

<!-- Update documentation for your changes -->

- [ ] Updated relevant documentation in `docs/`
- [ ] Updated code comments for exported functions
- [ ] Updated `CHANGELOG.md` (if applicable)
- [ ] Updated configuration examples (if applicable)

## Security

<!-- CRITICAL: Ensure no secrets are committed -->

- [ ] ✅ **NO SECRETS COMMITTED** (API keys, passwords, tokens, credentials)
- [ ] No hardcoded credentials in code
- [ ] Environment variables used for sensitive data
- [ ] Example configs use placeholders only
- [ ] Reviewed `git diff` for accidental secret inclusion

## Screenshots (if applicable)

<!-- Add screenshots for UI changes, terminal output, etc. -->

## Checklist

- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing tests pass locally with my changes
- [ ] I have checked my code doesn't commit any secrets

## Additional Notes

<!-- Any additional information that reviewers should know -->

---

**By submitting this PR, I confirm that**:
- I have read and followed the [Contributing Guidelines](../CONTRIBUTING.md)
- I understand the [Code of Conduct](../CODE_OF_CONDUCT.md)
- My code does not contain any secrets or sensitive information
