/**
 * buildWasm.js
 * 
 * This script builds the Go program in the repo base against wasm
 * and stores the output in website/public/SOMAS2020.wasm
 * 
 * You can run this script using `yarn buildwasm`.
 * This script also runs during prebuild and prestart.
 */


const path = require('path')
const cp = require('child_process')

try {
    const websiteRoot = path.dirname(__dirname)
    const repoRoot = path.dirname(websiteRoot)

    const stdout = cp.execSync(
        `go build -ldflags="-w -s" -o ./website/public/SOMAS2020.wasm`,
        {
            cwd: repoRoot,
            env: {
                ...process.env,
                GOOS: `js`,
                GOARCH: `wasm`,
            }
        }
    )
    console.log(stdout.toString())
}
catch (err) {
    console.error(err)
    process.exit(1)
}

console.log("Built SOMAS2020.wasm successfully")