# pluginlab

A small Go host for experimenting with plugin architectures: **compile-time registration + runtime configuration**, plus a subprocess plugin that runs external scripts.

The host builds a fake `BuildResult` and hands it to every notifier listed in the config.

## Run

```sh
go run .                # uses ./config.json
go run . other.json     # or point at another config
```

Exits non-zero if any notifier fails.

## How it works

- [types.go](types.go) — the `Notifier` interface (`Name` / `Configure` / `Notify`) and the `BuildResult` passed to plugins.
- [registry.go](registry.go) — `name -> Factory` map. Duplicate names panic at startup.
- Each `notify_*.go` calls `Register` in its `init()`, so linking the file in is what enables the plugin.
- [main.go](main.go) — reads the config, looks up each entry, `Configure`s it, then runs them in order. Unknown names warn and are skipped; a config error is fatal.

## Built-in notifiers

| name | options |
| --- | --- |
| `console` | `prefix` (default `[console]`) |
| `file` | `path` (required) — appends one line per build |
| `exec` | `command` (required), `args`, `transport` (`env` \| `json`, default `env`), `timeout` seconds (default 10) |

```json
{
  "notifiers": [
    { "name": "file", "config": { "path": "builds.log" } },
    { "name": "console", "config": { "prefix": "[ci]" } },
    { "name": "exec", "config": {
        "command": "./test_scripts/json_ok.sh",
        "transport": "json",
        "timeout": 1
    } }
  ]
}
```

## exec transports

The child gets a scrubbed env (`PATH` only), and is killed at `timeout`.

- **`env`** — build fields arrive as `BUILD_ID`, `BUILD_BRANCH`, `BUILD_STATUS`. stdout and stderr are both just log noise. Success = exit 0.
- **`json`** — the `BuildResult` is written to the child's stdin; the child must reply on stdout with `{"ok": bool, "message": string}`. **stdout is the protocol**, so log only to stderr. Success = exit 0 *and* `ok: true`.

## test_scripts

Fixtures for the failure modes worth seeing:

- `json_ok.sh` / `env_ok.sh` — happy path for each transport
- `fail.sh` — non-zero exit, stderr is surfaced
- `slow.sh` — `sleep 30`, trips the timeout
- `leak.sh` — dumps `env`; run it after deleting the `cmd.Env` assignment in [notify_exec.go](notify_exec.go) to see what the child would otherwise inherit

## Adding a notifier

Drop in a `notify_x.go` with a type implementing `Notifier` and an `init()` that calls `Register("x", ...)`. No other file changes.
