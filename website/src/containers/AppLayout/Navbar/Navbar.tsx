import React, { useState } from 'react'
import { Navbar, Nav, NavDropdown } from 'react-bootstrap';
import { LinkContainer } from 'react-router-bootstrap'
import { Link } from 'react-router-dom';
import logo from '../../../assets/logo/logo192.png';
import outputJSONData from '../../../output/output.json'

import styles from './Navbar.module.css'

const AppNavbar = () => {
  const [navExpanded, setNavExpanded] = useState(false)
  const closeNav = () => setNavExpanded(false)
  const getNavLink = (text: string, link: string) =>
    <LinkContainer to={link} onClick={closeNav}>
      <Nav.Link className="lightbluelink">{text}</Nav.Link>
    </LinkContainer>

  return <>
    <Navbar fixed="top" bg="dark" variant="dark" expand="lg"
      onToggle={() => setNavExpanded(!navExpanded)} expanded={navExpanded}>

      <Link to="/" className={styles.enlargeOnHover}>
        <Navbar.Brand>
          <img
            alt=""
            src={logo}
            width="30"
            height="30"
            className="d-inline-block align-top"
          />{' '}
                    SOMAS 2020
            </Navbar.Brand>
      </Link>

      <a rel="noopener noreferrer" target="_blank" href={outputJSONData.GitInfo.GithubURL} className="lightbluelink">
        {outputJSONData.GitInfo.Hash.substr(0, 7)}
      </a>

      <Navbar.Toggle aria-controls="responsive-navbar-nav" onClick={closeNav} />
      <Navbar.Collapse id="responsive-navbar-nav" className="justify-content-end">
        <Nav className="mr-auto" />
        <Nav >
          {getNavLink("Raw Output", "/rawoutput")}
          <NavDropdown title="Visualisations" id="collasible-nav-dropdown">
            <NavDropdown.Item href="/resources">Resources</NavDropdown.Item>
            <NavDropdown.Item href="/roles">Roles by Turn</NavDropdown.Item>
            <NavDropdown.Divider />
            <NavDropdown.Item href="/IIGO">IIGO</NavDropdown.Item>
            <NavDropdown.Item href="/IITO">IITO</NavDropdown.Item>
            <NavDropdown.Item href="/IIFO">IIFO</NavDropdown.Item>
          </NavDropdown>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  </>
}

export default AppNavbar