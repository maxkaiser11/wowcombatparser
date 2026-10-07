# WoW Combat Parser

A World of Warcraft combat log analyser written in Go. Upload a `WoWCombatLog.txt` file and get a damage and DPS breakdown for every boss encounter in it, from the command line or in the browser.

This is the first stage of **WoW Coach**, a tool that will eventually give players automated, AI-assisted coaching from their combat logs.

## Features

- **Streaming parser**: reads logs line by line, so multi-gigabyte raid logs never have to fit in memory
- **Encounter detection**: splits the log into boss pulls using `ENCOUNTER_START` / `ENCOUNTER_END`, with kill/wipe status and the official fight length
- **Per-player damage and DPS**: spell, periodic, ranged and melee damage, filtered to players only
- **Web interface**: upload a log in the browser and see a ranked DPS table per fight
- **CLI**: the same analysis in the terminal
- **Robust error handling**: malformed lines are counted and reported instead of crashing the run

## Project structure

```
.
├── cmd/
│   ├── server/          # web server entry point
│   └── wowparse/        # command-line tool entry point
└── internal/
    ├── api/             # HTTP handlers, routes and view models
    ├── combatlog/       # log parsing: lines, damage events, encounters
    └── web/
        ├── static/      # CSS
        └── templates/   # HTML templates (embedded into the binary)
```

The parser lives in `internal/combatlog` and has no knowledge of HTTP or the terminal. `combatlog.Parse` takes any `io.Reader`, so the CLI passes it a file, the server passes it an upload, and tests pass it a string.

## Getting started

Requires **Go 1.22+**.

```bash
git clone https://github.com/maxkaiser11/wowcombatparser.git
cd wowcombatparser
go mod tidy
```

### Recording a log

1. In WoW, enable **Advanced Combat Logging** (System → Network).
2. Type `/combatlog` before your raid or dungeon.
3. Logs are saved to `World of Warcraft/_retail_/Logs/`.

### Web server

```bash
go run ./cmd/server
```

Open [http://localhost:8080](http://localhost:8080) and upload a log.

### Command line

```bash
go run ./cmd/wowparse path/to/WoWCombatLog.txt
```

### Tests

```bash
go test ./...
```

## Endpoints

| Method | Path        | Description                           |
| ------ | ----------- | ------------------------------------- |
| GET    | `/`         | Upload page                           |
| POST   | `/upload`   | Parse an uploaded log, render results |
| GET    | `/health`   | Health check                          |
| GET    | `/static/*` | Static assets                         |

Uploads are sent as `multipart/form-data` in a field named `log`, with a limit of 1 GB.

## Known limitations

- Pet and guardian damage is not yet credited to their owners, so pet classes show lower DPS than in-game meters.
- Only damage is analysed; healing, deaths and cooldowns are planned.
- Logs are parsed during the request and not stored.

## Roadmap

- [x] Streaming combat log parser
- [x] Encounter splitting with kill/wipe and duration
- [x] Web upload and results page
- [ ] JSON API
- [ ] PostgreSQL storage and upload history
- [ ] Background parsing worker
- [ ] User accounts
- [ ] Docker, CI and deployment
- [ ] Death recaps, cooldown and buff uptime analysis
- [ ] AI coaching feedback based on detected mistakes