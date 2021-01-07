import React, { useEffect, useRef } from "react";
import { runForceGraph } from "./ForceGraphGenerator";
import styles from "../IITO.module.css";

const ForceGraph = ({ linksData, nodesData, nodeHoverTooltip }) => {
    const containerRef = useRef(null);

    useEffect(() => {
        let destroyFn;

        if (containerRef.current) {
            const { destroy } = runForceGraph(containerRef.current, linksData, nodesData, nodeHoverTooltip);
            destroyFn = destroy;
        }

        return destroyFn;
    }, []);

    return <div ref={containerRef} className={styles.container} />;
}

export default ForceGraph;
