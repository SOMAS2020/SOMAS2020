/**
 * copyOutput.js
 * 
 * This script copies the output from the completed Go program simulation
 * from the root of the repository.
 * 
 * You can run this script using `yarn copyoutput`.
 * This script also runs during prebuild and prestart.
 */

const fse = require('fs-extra')
const fs = require('fs')
const path = require('path')

try {
    const websiteRoot = path.dirname(path.dirname(__dirname))
    const repoRoot = path.dirname(websiteRoot)

    fse.copySync(
        path.join(repoRoot, `output`), 
        path.join(websiteRoot, `src`, `output`),
    )
    
    const logTxt = fs.readFileSync(
        path.join(websiteRoot, `src`, `output`, `log.txt`), 
        `utf8`)
    const logLines = logTxt.split(`\n`)
    
    fs.writeFileSync(path.join(websiteRoot, `src`, `output`, `log.txt.json`), JSON.stringify(logLines))
}
catch (err) {
    console.error(err)
    process.exit(1)
}

console.log("Output copied and processed successfully")