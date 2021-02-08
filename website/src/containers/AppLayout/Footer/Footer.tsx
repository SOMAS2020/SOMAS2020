import React from 'react'

import styles from './Footer.module.css'

const Footer = () => {
  return (
    <div className={styles.root}>
      <div className={styles.box}>
        <h5>Self-Organising Multi-Agent Systems 2020</h5>
        <p>Imperial College London EEE Department</p>
        <a
          rel="noopener noreferrer"
          target="_blank"
          href="https://github.com/SOMAS2020/SOMAS2020"
          className="lightbluelink"
        >
          GitHub
        </a>
      </div>
    </div>
  )
}

export default Footer
