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

### Parameters & Help
```bash
go run . --help
```

### Output
After running, the `output` directory will contain the output of the program.
- `output.json`: JSON file containing the game's historic states and configuration.
- `log.txt`: logs of the run

### Visualisation Website
See [`website/README.md`](website/README.md)

### WebAssembly Output

A script is provided to compile the program into WebAssembly for use in the website.

On Linux/maxOS,
```bash
./build_wasm.sh
```

On Windows,
```bash
build_cmd.cmd
```

## Testing
```bash
go test ./...
```

## Structure

### [`docs`](docs)
Important documents pertaining to codebase organisation, code conventions and project management. Read before writing code.

### [`internal`](internal)
Internal SOMAS2020 packages. Most development occurs here, including client and server code.

- [`clients`](internal/clients)
Individual team code goes into the respective folders in this directory.

- [`common`](internal/common)
Common utilities, or system-wide code such as game specification etc.

- [`server`](internal/server)
Self-explanatory.

### [`pkg`](pkg)
More generic packages dealing with general use-cases, such as system-related or file-operation utilities.

### [`website`](website)
Source code for visualisation website.
