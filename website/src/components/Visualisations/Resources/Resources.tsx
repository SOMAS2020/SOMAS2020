import React from "react";
import logo from "../../../assets/logo/logo512.png";
import styles from "./Resources.module.css";

const Resources = () => {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <p className={styles.text}>Resource Visualisation</p>
    </div>
  );
};

export default Resources;
