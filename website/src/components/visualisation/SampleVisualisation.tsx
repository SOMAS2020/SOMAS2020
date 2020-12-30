import React, { useEffect, useState } from 'react';
import drawChart from './BasicD3';

const SampleVis = () => {
    const [data, setData] = useState([]);

    useEffect(() => {
        drawChart(400, 600);
    }, []);
    return (
        <div>
            <h2>Sample Graph with D3 and React</h2>
            <div id="dr-chart"></div>
        </div>
    );
}

export default SampleVis;