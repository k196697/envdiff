# envdiff

Compare `.env` files across environments and highlight missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <file1> <file2>
```

**Example:**

```bash
envdiff .env.development .env.production
```

**Sample output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED values:
  - LOG_LEVEL: "debug" vs "info"
  - PORT: "3000" vs "8080"
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--keys-only` | Only compare key names, ignore values |
| `--quiet` | Exit with non-zero status if differences found (useful for CI) |
| `--format json` | Output results as JSON |

---

## Use in CI

```bash
envdiff --quiet .env.example .env.production || exit 1
```

---

## License

MIT © [yourusername](https://github.com/yourusername)