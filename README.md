# SOMAS 2020

Main repository for SOMAS2020 Coursework.

- [Setup & Rules](./docs/SETUP.md)
- [Infra Info](./docs/INFRA.md)

## Running code
```bash
# Approach 1
go run . # Linux and macOS: Use `sudo go run .` if you encounter any "Permission denied" errors.

# Approach 2
go build # build step
./SOMAS2020 # SOMAS2020.exe if you're on Windows. Use `sudo` on Linux and macOS as Approach 1 if required.
```

### Output
After running, the `output` directory will contain the output of the program.
- `output.json`: JSON file containing the game's historic states and configuration.
- `log.txt`: logs of the run

## Testing
```bash
go test ./...
```

## Structure

### `docs`
Important documents pertaining to codebase organisation, code conventions and project management. Read before writing code.

### `internal`
All code goes in here.

#### `clients`
Individual team code goes into the respective folders in this directory.

#### `common`
Common utilities, or system-wide code such as game specification etc.

#### `server`
Self-explanatory.
