# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

```bash
# Diff two versions of a secret
vaultdiff secret/myapp/config --v1 3 --v2 4

# Compare secrets across environments
vaultdiff --addr-a https://vault-staging:8200 --addr-b https://vault-prod:8200 \
  secret/myapp/config

# Audit all changes between versions
vaultdiff secret/myapp/config --from 1 --to latest --audit
```

Output highlights added, removed, and modified keys while masking sensitive values by default. Use `--unmask` to reveal full values (requires appropriate Vault policies).

---

## Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--unmask` | Show secret values in plain text | `false` |
| `--output` | Output format (`text`, `json`) | `text` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.9+ with KV v2 secrets engine

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)