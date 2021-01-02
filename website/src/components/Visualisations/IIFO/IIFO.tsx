import React from "react";
import logo from "../../../assets/logo/logo512.png";
import styles from "./IIFO.module.css";

const IIFO = () => {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <p className={styles.text}>IIFO Visualisation</p>
    </div>
  );
};

export default IIFO;
