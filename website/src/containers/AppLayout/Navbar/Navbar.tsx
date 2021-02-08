import React, { useState } from 'react'
import { Navbar, Nav } from 'react-bootstrap'
import { LinkContainer } from 'react-router-bootstrap'
import { about, cioutput, newrun, visualisations } from '../../../consts/paths'
import outputJSONData from '../../../output/output.json'

import logo from '../../../assets/logo/logo192.png'

import styles from './Navbar.module.css'

const AppNavbar = () => {
  const [navExpanded, setNavExpanded] = useState(false)
  const closeNav = () => setNavExpanded(false)
  const getNavLink = (text: string, link: string) => (
    <LinkContainer to={link} onClick={closeNav}>
      <Nav.Link className="lightbluelink">{text}</Nav.Link>
    </LinkContainer>
  )

  return (
    <>
      <Navbar
        fixed="top"
        bg="dark"
        variant="dark"
        expand="lg"
        onToggle={() => setNavExpanded(!navExpanded)}
        expanded={navExpanded}
      >
        {/* lhl2617: DO NOT WRAP BRAND WITH Link AS IT BREAKS OTHER COMPONENTS */}
        <Navbar.Brand
          href={process.env.PUBLIC_URL}
          className={styles.enlargeOnHover}
        >
          <img
            alt=""
            src={logo}
            width="30"
            height="30"
            className="d-inline-block align-top"
          />{' '}
          SOMAS 2020
        </Navbar.Brand>

        <a
          rel="noopener noreferrer"
          target="_blank"
          href={outputJSONData.GitInfo.GithubURL}
          className="lightbluelink"
        >
          {outputJSONData.GitInfo.Hash.substr(0, 7)}
        </a>

        <Navbar.Toggle aria-controls="basic-navbar-nav" onClick={closeNav} />
        <Navbar.Collapse id="basic-navbar-nav" className="justify-content-end">
          <Nav className="mr-auto" />
          <Nav>
            {getNavLink('About', about)}
            {getNavLink('New Run', newrun)}
            {getNavLink('CI Output', cioutput)}
            {getNavLink('Visualisations', visualisations)}
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    </>
  )
}

export default AppNavbar
