import React, { useState } from 'react'
import { Navbar, Nav } from 'react-bootstrap'
import { LinkContainer } from 'react-router-bootstrap'
import {
  gamevisualisation,
  iigovisualisation,
  iitovisualisation,
  iifovisualisation,
  rolesvisualisation,
  resourcesvisualisation,
  achievementsvisualisation,
  visualisations,
} from '../../consts/paths'

const VisualisationsNavbar = (props: { reset: () => any }) => {
  const [navExpanded, setNavExpanded] = useState(false)
  const closeNav = () => setNavExpanded(false)
  const getNavLink = (text: string, link: string) => (
    <LinkContainer to={link} onClick={closeNav}>
      <Nav.Link>{text}</Nav.Link>
    </LinkContainer>
  )
  const handleReset = () => {
    const { reset } = props
    reset()
    closeNav()
  }

  return (
    <>
      <Navbar
        bg="primary"
        variant="dark"
        expand="lg"
        onToggle={() => setNavExpanded(!navExpanded)}
        expanded={navExpanded}
      >
        <Navbar.Toggle aria-controls="basic-navbar-nav" onClick={closeNav} />
        <Navbar.Collapse id="basic-navbar-nav" className="justify-content-end">
          <Nav className="mr-auto">
            {getNavLink('Game', gamevisualisation)}
            {getNavLink('IIGO', iigovisualisation)}
            {getNavLink('IITO', iitovisualisation)}
            {getNavLink('IIFO', iifovisualisation)}
            {getNavLink('Roles', rolesvisualisation)}
            {getNavLink('Resources', resourcesvisualisation)}
            {getNavLink('Achievements', achievementsvisualisation)}
          </Nav>
          <Nav>
            <LinkContainer exact to={visualisations} onClick={handleReset}>
              <Nav.Link>Reset</Nav.Link>
            </LinkContainer>
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    </>
  )
}

export default VisualisationsNavbar
