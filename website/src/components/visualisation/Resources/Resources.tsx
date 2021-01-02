import React from 'react';
import styles from './Resources.module.css';
import LineRechartComponent from './LineGraph';

const Resources = () => {
  return (
    <div className={styles.root}>
      <h1>Resources over Time </h1>
      <p>Select teams to show/hide by cicking the team at the top of the chart. 
         Use the slider at the bottom to change which turns are displayed. </p>
      <LineRechartComponent />
    </div>
  );
}

export default Resources;
