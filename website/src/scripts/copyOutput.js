const fse = require('fs-extra')
const fs = require('fs')

try {
    fse.copySync(`../output`, `src/output`)
    
    const logTxt = fs.readFileSync(`src/output/log.txt`, `utf8`)
    const logLines = logTxt.split(`\n`)
    
    fs.writeFileSync(`src/output/log.txt.json`, JSON.stringify(logLines))
}
catch (err) {
    console.error(err)
    process.exit(1)
}

console.log("Output copied and processed successfully")