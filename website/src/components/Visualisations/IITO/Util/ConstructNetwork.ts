const fs = require('fs');

function constructNetwork(data) {
    let nodes = [];
    let links = [];

    let validIds = new Set(data.map(({ id }) => id))

    data.forEach(element => {
        // assign a paper with an id to a node
        const { id, title, paperAbstract, authors } = element;
        let node = { id, title, paperAbstract, authors }
        nodes.push(node);

        // map through incitations
        // console.log(node.inCitations);
        element.inCitations.forEach(inCitation => {
            let link = {};
            link.source = inCitation;
            link.target = element.id;

            if (validIds.has(link.source) && validIds.has(link.target)) {
                links.push(link)
            }
        })

        element.outCitations.forEach(outCitation => {
            let link = {};
            link.source = outCitation;
            link.target = element.id;

            if (validIds.has(link.source) && validIds.has(link.target)) {
                links.push(link)
            }
        })
    });

    // remove link source target pairs that reference nodes ids not in nodes

    let newArr = { nodes, links };

    // convert JSON object to string
    const newData = JSON.stringify(newArr);

    fs.writeFile('cleaned.json', newData, (err) => {
        if (err) {
            throw err;
        }
        console.log("JSON data is saved.");
    });
}

module.exports = {
    constructNetwork,
};