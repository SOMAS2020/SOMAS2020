import React from 'react'
import Content from './Content/Content'
import styles from './AppLayout.module.css'
import AppNavbar from './Navbar/Navbar'
import Footer from './Footer/Footer'

const AppLayout = () => {
  return (
    <div className={styles.root}>
      <AppNavbar />
      <div className={styles.content}>
        <Content />
      </div>
      <Footer />
    </div>
  )
}

export default AppLayout
