#!/usr/bin/env bash
#
# Release Automation Script for Termplate Go
#
# This script automates the release process by:
# - Validating prerequisites (git status, tests, lint)
# - Parsing CHANGELOG.md [Unreleased] section
# - Creating versioned section with current date
# - Updating version comparison links
# - Creating annotated git tag
# - Pushing tag to remote (triggers GitHub Actions)
#
# Usage:
#   ./scripts/release.sh                    # Interactive mode
#   ./scripts/release.sh --dry-run          # Preview changes
#   ./scripts/release.sh --version v0.2.0   # Specify version
#   ./scripts/release.sh --patch            # Auto-increment patch
#   ./scripts/release.sh --skip-tests       # Skip test validation
#

set -e  # Exit on error

# ============================================================================
# Configuration & Constants
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CHANGELOG_FILE="$PROJECT_ROOT/CHANGELOG.md"
CHANGELOG_BACKUP="$PROJECT_ROOT/CHANGELOG.md.backup"
REPO_OWNER="pranav3714"
REPO_NAME="termplate-go"

# ============================================================================
# Color Output Functions
# ============================================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1" >&2
}

header() {
    echo ""
    echo -e "${BOLD}$1${NC}"
    echo "$(printf '─%.0s' {1..60})"
}

# ============================================================================
# Global Variables
# ============================================================================

DRY_RUN=false
SKIP_TESTS=false
SKIP_LINT=false
AUTO_CONFIRM=false
VERSION=""
AUTO_INCREMENT=""
RELEASE_NOTES=""
PREV_VERSION=""  # Store previous version for comparison links

# ============================================================================
# Helper Functions
# ============================================================================

usage() {
    cat <<EOF
Release Automation Script for Termplate Go

Usage:
    $0 [OPTIONS]

Options:
    --version VERSION    Specify version (e.g., v0.2.0)
    --patch             Auto-increment patch version (0.1.1 → 0.1.2)
    --minor             Auto-increment minor version (0.1.1 → 0.2.0)
    --major             Auto-increment major version (0.1.1 → 1.0.0)
    --dry-run           Preview changes without committing
    --skip-tests        Skip running tests
    --skip-lint         Skip running linters
    --yes               Auto-confirm all prompts
    --help              Show this help message

Examples:
    $0                              # Interactive mode
    $0 --dry-run                    # Preview changes
    $0 --version v0.2.0             # Specify version
    $0 --patch                      # Auto-increment patch
    $0 --skip-tests --skip-lint     # Skip validation

EOF
    exit 0
}

# Parse version string (v1.2.3 → MAJOR=1 MINOR=2 PATCH=3)
parse_version() {
    local version="$1"
    if [[ $version =~ ^v?([0-9]+)\.([0-9]+)\.([0-9]+)(-.*)?$ ]]; then
        echo "${BASH_REMATCH[1]} ${BASH_REMATCH[2]} ${BASH_REMATCH[3]}"
        return 0
    else
        return 1
    fi
}

# Get current version from git tags
get_current_version() {
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Increment version based on type (patch/minor/major)
increment_version() {
    local current="$1"
    local type="$2"

    local parts
    parts=$(parse_version "$current") || {
        error "Failed to parse current version: $current"
        return 1
    }

    read -r major minor patch <<< "$parts"

    case "$type" in
        patch)
            patch=$((patch + 1))
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
    esac

    echo "v${major}.${minor}.${patch}"
}

# Validate version format (semver)
validate_version() {
    local version="$1"
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        error "Invalid version format: $version"
        error "Expected format: vMAJOR.MINOR.PATCH (e.g., v0.2.0 or v1.0.0-beta.1)"
        return 1
    fi
    return 0
}

# Confirm action with user
confirm() {
    if [[ "$AUTO_CONFIRM" == "true" ]]; then
        return 0
    fi

    local prompt="$1"
    local response
    read -p "$prompt [y/N] " -n 1 -r response
    echo
    [[ "$response" =~ ^[Yy]$ ]]
}

# ============================================================================
# Validation Functions
# ============================================================================

validate_prerequisites() {
    header "Validating Prerequisites"

    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        error "Not a git repository"
        return 1
    fi
    success "Git repository detected"

    # Check if working directory is clean
    if [[ -n $(git status --porcelain) ]]; then
        error "Working directory is not clean"
        info "Please commit or stash your changes first"
        echo ""
        git status --short
        return 1
    fi
    success "Working directory is clean"

    # Check if on main branch
    local current_branch
    current_branch=$(git branch --show-current)
    if [[ "$current_branch" != "main" ]]; then
        warning "Not on main branch (current: $current_branch)"
        if ! confirm "Continue anyway?"; then
            error "Aborted by user"
            return 1
        fi
    else
        success "On main branch"
    fi

    # Check if remote is configured
    if ! git remote get-url origin > /dev/null 2>&1; then
        error "No remote 'origin' configured"
        return 1
    fi
    success "Remote 'origin' configured"

    # Run tests unless skipped
    if [[ "$SKIP_TESTS" == "false" ]]; then
        info "Running tests..."
        local start_time=$(date +%s)
        if ! make test > /dev/null 2>&1; then
            error "Tests failed"
            info "Run 'make test' to see details"
            if ! confirm "Continue anyway?"; then
                error "Aborted by user"
                return 1
            fi
        else
            local end_time=$(date +%s)
            local duration=$((end_time - start_time))
            success "Tests passed (${duration}s)"
        fi
    else
        warning "Skipping tests (--skip-tests)"
    fi

    # Run linters unless skipped
    if [[ "$SKIP_LINT" == "false" ]]; then
        info "Running linters..."
        local start_time=$(date +%s)
        if ! make lint > /dev/null 2>&1; then
            error "Linters failed"
            info "Run 'make lint' to see issues"
            if ! confirm "Continue anyway?"; then
                error "Aborted by user"
                return 1
            fi
        else
            local end_time=$(date +%s)
            local duration=$((end_time - start_time))
            success "Linters passed (${duration}s)"
        fi
    else
        warning "Skipping linters (--skip-lint)"
    fi

    # Check if CHANGELOG.md exists
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        error "CHANGELOG.md not found"
        return 1
    fi
    success "CHANGELOG.md found"

    echo ""
    return 0
}

validate_version_available() {
    local version="$1"

    if git rev-parse "$version" >/dev/null 2>&1; then
        error "Tag $version already exists"
        info "Either delete the existing tag or choose a different version:"
        info "  git tag -d $version"
        info "  git push origin :refs/tags/$version  # Delete remote"
        return 1
    fi

    return 0
}

# ============================================================================
# CHANGELOG Functions
# ============================================================================

parse_unreleased_section() {
    local changelog="$1"
    local in_unreleased=false
    local unreleased_content=""
    local line_num=0

    while IFS= read -r line; do
        line_num=$((line_num + 1))

        # Start capturing after [Unreleased] header
        if [[ "$line" =~ ^\#\#[[:space:]]*\[Unreleased\] ]]; then
            in_unreleased=true
            continue
        fi

        # Stop at next ## heading
        if [[ "$in_unreleased" == "true" && "$line" =~ ^\#\#[[:space:]] ]]; then
            break
        fi

        # Capture content
        if [[ "$in_unreleased" == "true" ]]; then
            unreleased_content+="$line"$'\n'
        fi
    done < "$changelog"

    # Trim trailing newlines
    unreleased_content="${unreleased_content%$'\n'}"

    echo "$unreleased_content"
}

update_changelog() {
    local version="$1"
    local date="$2"
    local changelog="$CHANGELOG_FILE"

    header "Updating CHANGELOG.md"

    # Create backup
    cp "$changelog" "$CHANGELOG_BACKUP"
    info "Created backup: $CHANGELOG_BACKUP"

    # Store previous version before any modifications
    PREV_VERSION=$(grep -m 1 '^\[Unreleased\]:' "$changelog" | sed -n 's/.*compare\/\(v[0-9.]*\)\.\.\.HEAD/\1/p')
    if [[ -z "$PREV_VERSION" ]]; then
        PREV_VERSION=$(get_current_version)
    fi
    info "Previous version: $PREV_VERSION"

    # Parse unreleased section
    local unreleased_content
    unreleased_content=$(parse_unreleased_section "$changelog")

    # Check if unreleased section is empty
    if [[ -z "$(echo "$unreleased_content" | tr -d '[:space:]')" ]]; then
        warning "[Unreleased] section is empty"
        if [[ "$AUTO_CONFIRM" == "false" ]]; then
            warning "No changes documented for this release"
            if ! confirm "Continue anyway?"; then
                error "Aborted by user"
                return 1
            fi
        fi
    else
        echo ""
        info "Changes in this release:"
        echo "$unreleased_content" | head -15
        if [[ $(echo "$unreleased_content" | wc -l) -gt 15 ]]; then
            info "... ($(echo "$unreleased_content" | wc -l) total lines)"
        fi
        echo ""
    fi

    # Read entire changelog
    local changelog_content
    changelog_content=$(cat "$changelog")

    # Create new version section
    local version_section="## [${version#v}] - $date"

    # Build new changelog structure
    local new_changelog=""
    local found_unreleased=false
    local in_unreleased=false
    local skip_unreleased_content=false
    local line_num=0

    while IFS= read -r line; do
        line_num=$((line_num + 1))

        # Found [Unreleased] heading
        if [[ "$line" =~ ^\#\#[[:space:]]*\[Unreleased\] ]]; then
            found_unreleased=true
            in_unreleased=true
            skip_unreleased_content=true
            # Write header up to and including [Unreleased]
            new_changelog+="$line"$'\n'$'\n'
            continue
        fi

        # Skip the old unreleased content (we'll replace it)
        if [[ "$skip_unreleased_content" == "true" && "$in_unreleased" == "true" ]]; then
            # Check if we've reached the next ## heading (end of unreleased section)
            if [[ "$line" =~ ^\#\#[[:space:]] ]]; then
                in_unreleased=false
                skip_unreleased_content=false
                # Insert new version section and then continue with the old version
                new_changelog+="$version_section"$'\n'
                new_changelog+="$unreleased_content"$'\n'$'\n'
                new_changelog+="$line"$'\n'
            fi
            continue
        fi

        # Append everything else
        if [[ "$skip_unreleased_content" == "false" ]]; then
            new_changelog+="$line"$'\n'
        fi
    done < "$changelog"

    if [[ "$found_unreleased" == "false" ]]; then
        error "Could not find [Unreleased] section in CHANGELOG.md"
        error "Please ensure CHANGELOG.md follows Keep a Changelog format"
        rm "$CHANGELOG_BACKUP"
        return 1
    fi

    # Write updated changelog
    if [[ "$DRY_RUN" == "false" ]]; then
        echo "$new_changelog" > "$changelog"
        success "Updated CHANGELOG.md"
    else
        info "[DRY RUN] Would update CHANGELOG.md with:"
        echo ""
        echo "  ## [${version#v}] - $date"
        echo "$unreleased_content" | head -10 | sed 's/^/  /'
        if [[ $(echo "$unreleased_content" | wc -l) -gt 10 ]]; then
            info "  ... (and $(( $(echo "$unreleased_content" | wc -l) - 10 )) more lines)"
        fi
        echo ""
    fi

    return 0
}

update_version_links() {
    local version="$1"
    local changelog="$CHANGELOG_FILE"

    info "Updating version comparison links..."

    # Use the PREV_VERSION that was stored before CHANGELOG modification
    local prev_version="$PREV_VERSION"

    if [[ -z "$prev_version" ]]; then
        error "Previous version not detected"
        prev_version=$(get_current_version)
        warning "Falling back to git tag: $prev_version"
    fi

    # Read current links section
    local links_start
    links_start=$(grep -n '^\[Unreleased\]:' "$changelog" | cut -d: -f1)

    if [[ -z "$links_start" ]]; then
        error "Could not find [Unreleased] link in CHANGELOG.md"
        return 1
    fi

    # Create new links
    local new_unreleased_link="[Unreleased]: https://github.com/$REPO_OWNER/$REPO_NAME/compare/$version...HEAD"
    local new_version_link="[${version#v}]: https://github.com/$REPO_OWNER/$REPO_NAME/compare/$prev_version...$version"

    # Update links
    if [[ "$DRY_RUN" == "false" ]]; then
        # Replace [Unreleased] link
        sed -i "${links_start}s|.*|$new_unreleased_link|" "$changelog"
        # Insert new version link after [Unreleased]
        sed -i "${links_start}a\\$new_version_link" "$changelog"
        success "Updated version links"
        info "  [Unreleased]: ...compare/$version...HEAD"
        info "  [${version#v}]: ...compare/$prev_version...$version"
    else
        info "[DRY RUN] Would update links:"
        info "  [Unreleased]: ...compare/$version...HEAD"
        info "  [${version#v}]: ...compare/$prev_version...$version"
    fi

    return 0
}

# ============================================================================
# Git Functions
# ============================================================================

create_git_tag() {
    local version="$1"

    header "Creating Git Tag"

    # Stage CHANGELOG.md
    if [[ "$DRY_RUN" == "false" ]]; then
        git add "$CHANGELOG_FILE"
        success "Staged CHANGELOG.md"
    else
        info "[DRY RUN] Would stage: CHANGELOG.md"
    fi

    # Create commit
    local commit_message="chore: prepare release $version"
    if [[ "$DRY_RUN" == "false" ]]; then
        git commit -m "$commit_message"$'\n\n'"Co-Authored-By: Claude Sonnet 4.5 (1M context) <noreply@anthropic.com>"
        local commit_sha=$(git rev-parse --short HEAD)
        success "Created commit $commit_sha"
        info "Message: $commit_message"
    else
        info "[DRY RUN] Would create commit:"
        info "  $commit_message"
    fi

    # Get release notes from CHANGELOG
    local release_notes
    release_notes=$(parse_unreleased_section "$CHANGELOG_BACKUP")

    # Create annotated tag
    local tag_message="Release $version"$'\n\n'"$release_notes"$'\n\n'"Co-Authored-By: Claude Sonnet 4.5 (1M context) <noreply@anthropic.com>"

    if [[ "$DRY_RUN" == "false" ]]; then
        git tag -a "$version" -m "$tag_message"
        success "Created annotated tag $version"
        echo ""
        info "Tag contains:"
        echo "$release_notes" | head -8 | sed 's/^/  /'
        if [[ $(echo "$release_notes" | wc -l) -gt 8 ]]; then
            info "  ... ($(echo "$release_notes" | wc -l) total lines)"
        fi
    else
        info "[DRY RUN] Would create annotated tag: $version"
        echo ""
        info "Tag would contain:"
        echo "$release_notes" | head -8 | sed 's/^/  /'
        if [[ $(echo "$release_notes" | wc -l) -gt 8 ]]; then
            info "  ... ($(echo "$release_notes" | wc -l) total lines)"
        fi
    fi

    return 0
}

push_tag() {
    local version="$1"

    header "Pushing to Remote"

    if [[ "$DRY_RUN" == "true" ]]; then
        info "[DRY RUN] Would push tag to remote"
        info "Commands:"
        info "  git push origin main"
        info "  git push origin $version"
        return 0
    fi

    info "This will push the tag to remote and trigger GitHub Actions"
    if ! confirm "Push tag $version to remote?"; then
        warning "Tag created locally but not pushed"
        info "You can push it later with:"
        info "  git push origin main"
        info "  git push origin $version"
        return 0
    fi

    info "Pushing commit to main branch..."
    git push origin main
    success "Pushed commit"

    info "Pushing tag $version..."
    git push origin "$version"
    success "Pushed tag"

    echo ""
    success "Release $version initiated!"
    info "GitHub Actions triggered: https://github.com/$REPO_OWNER/$REPO_NAME/actions"
    echo ""

    return 0
}

# ============================================================================
# Interactive Functions
# ============================================================================

prompt_version() {
    if [[ -n "$VERSION" ]]; then
        return 0
    fi

    header "Version Selection"

    local current_version
    current_version=$(get_current_version)
    info "Current version: $current_version"

    # Calculate suggested versions
    local patch minor major
    patch=$(increment_version "$current_version" "patch")
    minor=$(increment_version "$current_version" "minor")
    major=$(increment_version "$current_version" "major")

    echo ""
    echo "Suggested versions:"
    echo "  1) $patch (patch - bug fixes)"
    echo "  2) $minor (minor - new features)"
    echo "  3) $major (major - breaking changes)"
    echo "  4) Custom version"
    echo ""

    local choice
    read -p "Select version [1-4]: " choice

    case "$choice" in
        1) VERSION="$patch" ;;
        2) VERSION="$minor" ;;
        3) VERSION="$major" ;;
        4)
            read -p "Enter version (e.g., v0.2.0): " VERSION
            ;;
        *)
            error "Invalid choice"
            return 1
            ;;
    esac

    validate_version "$VERSION" || return 1
    validate_version_available "$VERSION" || return 1

    success "Selected version: $VERSION"
    return 0
}

show_summary() {
    local version="$1"

    header "Release Summary"

    echo "Version:  $version"
    echo "Date:     $(date +%Y-%m-%d)"
    echo "Dry run:  $DRY_RUN"
    echo ""

    if [[ "$DRY_RUN" == "true" ]]; then
        warning "DRY RUN MODE - No changes will be made"
        info "Run without --dry-run to actually perform the release"
    fi
}

# ============================================================================
# Main Function
# ============================================================================

main() {
    cd "$PROJECT_ROOT"

    header "Termplate Go Release Automation"

    # Validate prerequisites
    validate_prerequisites || exit 1

    # Prompt for version if not specified
    prompt_version || exit 1

    # Show summary
    show_summary "$VERSION"

    # Confirm before proceeding
    if [[ "$DRY_RUN" == "false" ]]; then
        echo ""
        if ! confirm "Proceed with release $VERSION?"; then
            error "Aborted by user"
            exit 1
        fi
    fi

    # Update CHANGELOG.md
    local date
    date=$(date +%Y-%m-%d)
    update_changelog "$VERSION" "$date" || exit 1

    # Update version links
    update_version_links "$VERSION" || exit 1

    # Create git tag
    create_git_tag "$VERSION" || exit 1

    # Push to remote
    push_tag "$VERSION" || exit 1

    # Cleanup backup on success
    if [[ "$DRY_RUN" == "false" && -f "$CHANGELOG_BACKUP" ]]; then
        rm "$CHANGELOG_BACKUP"
    fi

    header "Done!"

    if [[ "$DRY_RUN" == "true" ]]; then
        echo ""
        info "This was a dry run. No changes were made."
        info "Run without --dry-run to perform the actual release:"
        echo ""
        echo "  ./scripts/release.sh --version $VERSION"
        echo "  # or"
        echo "  make release-prepare"
        echo ""
    else
        echo ""
        success "Release $VERSION completed successfully!"
        echo ""
        info "What happens next:"
        info "  1. GitHub Actions builds binaries (6 platforms)"
        info "  2. Release created: https://github.com/$REPO_OWNER/$REPO_NAME/releases/tag/$VERSION"
        info "  3. Monitor progress: https://github.com/$REPO_OWNER/$REPO_NAME/actions"
        echo ""
        info "Files modified:"
        info "  • CHANGELOG.md (version $VERSION added)"
        info "  • Git commit: chore: prepare release $VERSION"
        info "  • Git tag: $VERSION"
        echo ""
        warning "Rollback instructions (if needed):"
        echo "  git tag -d $VERSION                    # Delete local tag"
        echo "  git push origin :refs/tags/$VERSION    # Delete remote tag"
        echo "  git reset HEAD~1                       # Undo commit"
        echo "  git checkout CHANGELOG.md              # Restore CHANGELOG"
        echo ""
    fi
}

# ============================================================================
# Argument Parsing
# ============================================================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --patch)
            AUTO_INCREMENT="patch"
            shift
            ;;
        --minor)
            AUTO_INCREMENT="minor"
            shift
            ;;
        --major)
            AUTO_INCREMENT="major"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --skip-lint)
            SKIP_LINT=true
            shift
            ;;
        --yes)
            AUTO_CONFIRM=true
            shift
            ;;
        --help)
            usage
            ;;
        *)
            error "Unknown option: $1"
            usage
            ;;
    esac
done

# Handle auto-increment
if [[ -n "$AUTO_INCREMENT" ]]; then
    current_version=$(get_current_version)
    VERSION=$(increment_version "$current_version" "$AUTO_INCREMENT")
    info "Auto-incrementing $AUTO_INCREMENT: $current_version → $VERSION"
fi

# Run main function
main
