import React, { useState } from 'react'
import { Navbar, Nav, NavDropdown } from 'react-bootstrap';
import { LinkContainer } from 'react-router-bootstrap'
import {
  cioutput,
  newrun,
  gamevisualisation,
  iigovisualisation,
  iitovisualisation,
  iifovisualisation,
  resourcesvisualisation,
  rolesvisualisation
} from '../../../consts/paths';

import logo from '../../../assets/logo/logo192.png';

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

      <Navbar.Brand href="/" className={styles.enlargeOnHover}>
        <img
          alt=""
          src={logo}
          width="30"
          height="30"
          className="d-inline-block align-top"
        />{' '}
                    SOMAS 2020
            </Navbar.Brand>
      <Navbar.Toggle aria-controls="basic-navbar-nav" onClick={closeNav} />
      <Navbar.Collapse id="basic-navbar-nav" className="justify-content-end">
        <Nav className="mr-auto" />
        <Nav>
          {getNavLink("New Run", newrun)}
          {getNavLink("CI Output", cioutput)}
          <NavDropdown title="Visualisations" id="collabsible-nav-dropdown">
            <NavDropdown.Item href={gamevisualisation}>Game</NavDropdown.Item>
            <NavDropdown.Item href={iigovisualisation}>IIGO</NavDropdown.Item>
            <NavDropdown.Item href={iitovisualisation}>IITO</NavDropdown.Item>
            <NavDropdown.Item href={iifovisualisation}>IIFO</NavDropdown.Item>
            <NavDropdown.Item href={resourcesvisualisation}>Resources</NavDropdown.Item>
            <NavDropdown.Item href={rolesvisualisation}>Roles</NavDropdown.Item>
          </NavDropdown>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  </>
}

export default AppNavbar