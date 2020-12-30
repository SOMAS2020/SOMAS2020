import React, { useRef, useEffect, useState } from 'react';
import * as d3 from 'd3';

interface IProps {
    data?: any[];
}

const SampleVis = (props: IProps) => {
    // useRef Hook creates a variable that keeps values across rendering passes 
    // holds the component's SVG DOM element
    // initialised as null and React assigns it later (see return statement)
    const d3Container = useRef(null);

    useEffect(() => {
        if (props.data && d3Container.current) {
            const svg = d3.select(d3Container.current);

            // Bind D3 data
            const update = svg
                .append('g')
                .selectAll('text')
                .data(props.data);

            // Enter new D3 elements
            update.enter()
                .append('text')
                .attr('x', (d, i) => i * 25)
                .attr('y', 40)
                .style('font-size', 24)
                .text((d: number) => d);

            // Update existing D3 elements
            update
                .attr('x', (d, i) => i * 40)
                .text((d: number) => d);

            // Remove old D3 elements
            update.exit()
                .remove();
        }
    },

        /*
            useEffect has a dependency array (below). It's a list of dependency
            variables for this useEffect block. The block will run after mount
            and whenever any of these variables change. We still have to check
            if the variables are valid, but we do not have to compare old props
            to next props to decide whether to rerender.
        */
        [props.data, d3Container.current])

    return (
        <svg
            className="d3-component"
            width={400}
            height={200}
            ref={d3Container}
        />
    );
}

export default SampleVis;