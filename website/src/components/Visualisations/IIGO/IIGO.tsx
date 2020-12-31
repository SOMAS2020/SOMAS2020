import React from 'react';
import logo from '../../../assets/logo/logo512.png';
import styles from './IIGO.module.css';

const IIGO = () => {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <p className={styles.text}>IIGO Visualisation</p>
    </div>
  );
}

export default IIGO;
