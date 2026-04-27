# envlock

> A CLI tool to snapshot, diff, and restore environment variable sets across dev/staging/prod configs.

---

## Installation

```bash
go install github.com/yourusername/envlock@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlock.git
cd envlock
go build -o envlock .
```

---

## Usage

**Snapshot** your current environment:

```bash
envlock snapshot --env dev --output dev.lock
```

**Diff** two environment snapshots:

```bash
envlock diff dev.lock staging.lock
```

**Restore** an environment from a snapshot:

```bash
envlock restore --file dev.lock
```

**Example output:**

```
+ API_URL=https://staging.example.com
- API_URL=https://dev.example.com
~ LOG_LEVEL: debug → info
```

---

## Commands

| Command    | Description                              |
|------------|------------------------------------------|
| `snapshot` | Capture current environment variables   |
| `diff`     | Compare two environment snapshots        |
| `restore`  | Apply a saved snapshot to the shell      |

---

## License

[MIT](LICENSE) © 2024 yourusername