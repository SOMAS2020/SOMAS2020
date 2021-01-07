import React from 'react';
import logo from '../../../assets/logo/logo512.png';
import styles from './Game.module.css';
import { OutputJSONType } from "../../../consts/types";
// import SampleVis from '../visualisation/SampleVis';

const Game = (props: { output: OutputJSONType }) => {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <h1>Game Visualisation</h1>
      {/* <SampleVis data={[1, 2, 3]} /> */}
    </div>
  );
}

export default Game;
