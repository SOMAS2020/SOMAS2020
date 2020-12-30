import React from 'react';
import logo from '../../assets/logo/logo512.png';
import styles from './Home.module.css';
import Game from '../Game/Game';

function Home() {
  return (
    <div className={styles.root}>
      <Game />
    </div>
  );
}

export default Home;
