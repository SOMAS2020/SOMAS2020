/**
 * copyOutput.js
 *
 * This script copies the output from the completed Go program simulation
 * from the root of the repository.
 *
 * The logs are processed into a JSON list for ease of inclusion in the
 * React App.
 *
 * If the output is not found, this scripts attempts to run the simulation.
 *
 * You can run this script using `yarn copyoutput`.
 * This script also runs during prebuild and prestart.
 */

const fse = require('fs-extra')
const fs = require('fs')
const path = require('path')
const cp = require('child_process')

try {
  const websiteRoot = path.dirname(__dirname)
  const repoRoot = path.dirname(websiteRoot)

  if (!fs.existsSync(path.join(repoRoot, `output`))) {
    console.log('Did not find simulation output, running simulation...')
    const stdout = cp.execSync(`go run .`, {
      cwd: repoRoot,
    })
    console.log(stdout.toString())
  }

  fse.copySync(
    path.join(repoRoot, `output`),
    path.join(websiteRoot, `src`, `output`)
  )

  const logTxt = fs.readFileSync(
    path.join(websiteRoot, `src`, `output`, `log.txt`),
    `utf8`
  )
  const logLines = logTxt.split(`\n`)

  fs.writeFileSync(
    path.join(websiteRoot, `src`, `output`, `log.txt.json`),
    JSON.stringify(logLines)
  )
} catch (err) {
  console.error(err)
  process.exit(1)
}

console.log('Output copied and processed successfully')
