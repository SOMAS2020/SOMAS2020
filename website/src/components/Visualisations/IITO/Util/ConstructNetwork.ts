const fs = require('fs');

// Construct a network of nodes and links from the processed Transaction Data
function constructNetwork(data) {
    let nodes = [];
    let links = [];

    let validIds = new Set(data.map(({ id }) => id))

    data.forEach(element => {
        // assign a paper with an id to a node
        const { id, title, paperAbstract, authors } = element;
        let node = { id, title, paperAbstract, authors }
        nodes.push(node);

        // map through incitations -> will be transactions in
        element.inCitations.forEach(inCitation => {
            let link = {};
            link.source = inCitation;
            link.target = element.id;

            if (validIds.has(link.source) && validIds.has(link.target)) {
                links.push(link)
            }
        })

        // map through incitations -> will be transactions out
        element.outCitations.forEach(outCitation => {
            let link = {};
            link.source = outCitation;
            link.target = element.id;

            if (validIds.has(link.source) && validIds.has(link.target)) {
                links.push(link)
            }
        })
    });

    // TODO: remove link source target pairs that reference nodes ids not in nodes
    let newArr = { nodes, links };

    // convert JSON object to string
    const networkData = JSON.stringify(newArr);

    return networkData;
}

module.exports = {
    constructNetwork,
};