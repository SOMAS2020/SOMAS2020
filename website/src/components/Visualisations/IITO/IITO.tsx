import React from "react";
import logo from "../../../assets/logo/logo512.png";
import styles from "./IITO.module.css";
import { OutputJSONType } from "../../../consts/types";
import { ForceGraph } from "./Util/ForceGraph";

const IITO = (props: { output: OutputJSONType }) => {
  const nodeHoverTooltip = React.useCallback((node) => {
    return `<div>${node.name}</div>`;
  }, []);

  return (
    <div className={styles.root}>
      <div style={{ textAlign: 'center' }}>
        <div style={{ marginTop: '50px' }}>
          <ForceGraph linksData={data.links} nodesData={data.nodes} nodeHoverTooltip={nodeHoverTooltip} />
        </div>
      </div>
    </div>
  );
};

export default IITO;
