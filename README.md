# gpc - Google Play Console CLI

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**gpc** is a fast, lightweight, and scriptable CLI for Google Play Console. Inspired by [App Store Connect CLI](https://github.com/rudrankriyam/App-Store-Connect-CLI), it provides comprehensive automation for Android app publishing workflows.

## Features

- 📦 **Release Management** - Upload bundles/APKs, manage tracks, staged rollouts
- 🧪 **Testing** - Internal sharing, tester management, track configuration
- 📝 **Store Listings** - Localized metadata, screenshots, graphics
- ⭐ **Reviews** - List, filter, and reply to user reviews
- 💰 **Monetization** - In-app products, subscriptions, purchase verification
- 👥 **Access Control** - User permissions and grants
- 🔧 **CI/CD Ready** - JSON output, environment variables, clean exit codes

## Design Philosophy

- **JSON-first output** - Machine-readable by default
- **Explicit over cryptic** - Clear flag names, no magic
- **No interactive prompts** - All inputs via flags or environment
- **Clean exit codes** - 0=success, 1=error, 2=validation

## Installation

### Install Script (Recommended)

**Public repo:**
```bash
curl -fsSL https://raw.githubusercontent.com/anthropics/gpc/main/install.sh | bash
```

**Private repo (requires GitHub token):**
```bash
export GITHUB_TOKEN="your_github_token"
curl -fsSL -H "Authorization: token $GITHUB_TOKEN" \
  https://raw.githubusercontent.com/YOUR_ORG/gpc/main/install.sh | \
  GPC_REPO="YOUR_ORG/gpc" GITHUB_TOKEN=$GITHUB_TOKEN bash
```

### Homebrew (macOS/Linux)

```bash
brew tap anthropics/tap
brew install gpc
```

### Go Install

```bash
go install github.com/anthropics/gpc/cmd/gpc@latest
```

### From Source

```bash
git clone https://github.com/anthropics/gpc.git
cd gpc
make install
```

### Download Binary

Download from [Releases](https://github.com/anthropics/gpc/releases).

## Quick Start

### 1. Create a Service Account

1. Go to [Google Cloud Console](https://console.cloud.google.com/iam-admin/serviceaccounts)
2. Create a new service account
3. Download the JSON key file
4. Enable the Google Play Developer API in your project

### 2. Grant Access in Play Console

1. Go to [Play Console API Access](https://play.google.com/console/developers/api-access)
2. Link your Google Cloud project
3. Grant access to your service account

### 3. Configure gpc

```bash
# Login with service account
gpc auth login --name "default" --credentials /path/to/service-account.json

# Or use environment variables
export GPC_CREDENTIALS_PATH=/path/to/service-account.json
export GPC_PACKAGE=com.example.app
```

### 4. Start Using

```bash
# List tracks
gpc tracks list --package com.example.app

# Upload a bundle
gpc bundles upload --package com.example.app --file app.aab --track internal

# Check reviews
gpc reviews list --package com.example.app --min-rating 1 --max-rating 3
```

## Commands

### Authentication

```bash
gpc auth login --name "profile" --credentials /path/to/creds.json
gpc auth switch --name "profile"
gpc auth list
gpc auth current
```

### Apps

```bash
gpc apps get --package com.example.app
gpc apps data-safety get --package com.example.app
```

### Tracks (Releases)

```bash
gpc tracks list --package com.example.app
gpc tracks get --package com.example.app --track production
gpc tracks update --package com.example.app --track production \
  --version-code 42 --rollout-percentage 10
gpc tracks promote --from internal --to beta --version-code 42
gpc tracks halt --track production
gpc tracks complete --track production
```

### Bundles & APKs

```bash
# Upload bundle (recommended)
gpc bundles upload --package com.example.app --file app.aab --track internal

# Upload APK (legacy)
gpc apks upload --package com.example.app --file app.apk --track internal

# List builds
gpc bundles list --package com.example.app
```

### Testing

```bash
# Internal sharing (instant distribution)
gpc testing internal-sharing upload --package com.example.app --file app.aab

# Manage testers
gpc testing testers list --track alpha
gpc testing testers add --track alpha --emails "a@test.com,b@test.com"
gpc testing testers remove --track alpha --emails "old@test.com"
```

### Store Listings

```bash
# List localizations
gpc listings list --package com.example.app

# Update listing
gpc listings update --package com.example.app --locale en-US \
  --title "My App" --short-description "A great app"

# Sync from directory (fastlane-compatible)
gpc listings sync --dir ./metadata/
```

### Screenshots & Graphics

```bash
# List images
gpc images list --locale en-US --type phoneScreenshots

# Upload
gpc images upload --locale en-US --type phoneScreenshots --file screenshot.png

# Sync from directory
gpc images sync --dir ./screenshots/
```

### Reviews

```bash
# List reviews
gpc reviews list --package com.example.app --min-rating 1 --max-rating 3

# Reply to review
gpc reviews reply --review-id "abc123" --text "Thank you for your feedback!"
```

### In-App Products

```bash
gpc products list --package com.example.app
gpc products get --sku premium_upgrade
gpc products create --sku new_item --title "New Item" --price-usd 0.99
gpc products delete --sku old_item --confirm
```

### Subscriptions

```bash
gpc subscriptions list --package com.example.app
gpc subscriptions get --product-id monthly_sub
gpc subscriptions base-plans list --product-id monthly_sub
```

### Purchases

```bash
# Verify purchase
gpc purchases verify --token "purchase_token" --product-id "item_sku"

# Check subscription status
gpc purchases subscription-status --token "sub_token" --product-id "sub_id"

# List voided purchases
gpc purchases voided list --start-time "2024-01-01T00:00:00Z"
```

### Users & Permissions

```bash
gpc users list
gpc users grant --email "dev@company.com" --role releaseManager
gpc users revoke --email "old@company.com" --confirm
```

### Edits (Advanced)

```bash
# Manual edit workflow
gpc edits create
gpc edits validate --edit-id "abc123"
gpc edits commit --edit-id "abc123"
gpc edits delete --edit-id "abc123" --confirm
```

## Output Formats

```bash
# JSON (default)
gpc tracks list

# Pretty JSON
gpc tracks list --pretty

# Table
gpc tracks list --output table

# TSV (for scripting)
gpc tracks list --output tsv

# Minimal (values only)
gpc tracks list --output minimal
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GPC_CREDENTIALS_PATH` | Path to service account JSON |
| `GPC_CREDENTIALS_B64` | Base64-encoded credentials |
| `GPC_PACKAGE` | Default package name |
| `GPC_PROFILE` | Default auth profile |
| `GPC_OUTPUT` | Default output format |
| `GPC_TIMEOUT` | Request timeout |
| `GPC_DEBUG` | Enable debug logging |

## CI/CD Integration

### GitHub Actions

```yaml
- name: Upload to Play Store
  env:
    GPC_CREDENTIALS_B64: ${{ secrets.PLAY_STORE_CREDENTIALS }}
  run: |
    gpc bundles upload \
      --package com.example.app \
      --file app/build/outputs/bundle/release/app-release.aab \
      --track internal
```

### Fastlane Migration

Replace fastlane supply with gpc:

```bash
# Instead of: fastlane supply
gpc listings sync --dir ./fastlane/metadata/android/
gpc images sync --dir ./fastlane/screenshots/android/
gpc bundles upload --file app.aab --track production
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error |
| 2 | Validation failure |

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Inspired by [App Store Connect CLI](https://github.com/rudrankriyam/App-Store-Connect-CLI) by Rudrank Riyam.
