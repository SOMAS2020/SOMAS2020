import React from "react";
import logo from "../../../assets/logo/logo512.png";
import styles from "./Roles.module.css";
import { OutputJSONType } from "../../../consts/types";

const Roles = (props: { output: OutputJSONType }) => {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <p className={styles.text}>Role Visualisation</p>
    </div>
  );
};

export default Roles;
