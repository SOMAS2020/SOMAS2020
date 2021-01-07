import React from "react";
import logo from "../../../assets/logo/logo512.png";
import styles from "./IITO.module.css";
import { OutputJSONType } from "../../../consts/types";

import processTransactionData from "./Util/ProcessTransactionData"
import { ForceGraph } from "./Util/ForceGraph";

const IITO = (props: { output: OutputJSONType }) => {
  const nodeHoverTooltip = React.useCallback((node) => {
    return `<div>${node.name}</div>`;
  }, []);

  function processOutput(output: OutputJSONType) {
    const transactionData = processTransactionData(output)
    return transactionData;
  }

  const data = processOutput(props.output);

  return (
    <div className={styles.root}>
      <div style={{ textAlign: 'center' }}>
        <ForceGraph linksData={data.links} nodesData={data.nodes} nodeHoverTooltip={nodeHoverTooltip} />
      </div>
    </div>
  );
};

export default IITO;
