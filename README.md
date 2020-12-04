# SOMAS 2020

Main repository for SOMAS2020 Coursework.

- [Setup & Rules](./docs/SETUP.md)
- [Infra Info](./docs/INFRA.md)

## Running code
```bash
# Approach 1
go run .

# Approach 2
go build # build step
./SOMAS2020 # SOMAS2020.exe if you're on Windows
```

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

test
