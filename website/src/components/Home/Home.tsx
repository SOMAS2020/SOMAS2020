import React from 'react';
import logo from '../../assets/logo/logo512.png';
import styles from './Home.module.css';

function Home() {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <p className={styles.text}>
        Self-Organising Multi-Agent Systems 2020
        </p>
    </div>
  );
}

export default Home;
